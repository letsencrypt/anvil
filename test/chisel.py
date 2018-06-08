"""
A simple client that uses the Python ACME library to run a test issuance against
a local Boulder server. Usage:

$ virtualenv venv
$ . venv/bin/activate
$ pip install -r requirements.txt
$ python chisel.py foo.com bar.com
"""
import json
import logging
import os
import socket
import sys
import threading
import time
import urllib2

from cryptography import x509
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.primitives.serialization import load_pem_private_key

import OpenSSL
from OpenSSL import SSL
import josepy

from acme import challenges
from acme import client as acme_client
from acme import errors as acme_errors
from acme import messages
from acme import standalone

logger = logging.getLogger()
logging.basicConfig()
logger.setLevel(int(os.getenv('LOGLEVEL', 20)))

DIRECTORY = os.getenv('DIRECTORY', 'http://localhost:4000/directory')

# URLs for management interface of challtestsrv
SET_TXT = "http://localhost:8055/set-txt"
CLEAR_TXT = "http://localhost:8055/clear-txt"

os.environ.setdefault('REQUESTS_CA_BUNDLE', 'test/wfe-tls/minica.pem')

def make_client(email=None):
    """Build an acme.Client and register a new account with a random key."""
    key = josepy.JWKRSA(key=rsa.generate_private_key(65537, 2048, default_backend()))

    net = acme_client.ClientNetwork(key, user_agent="Boulder integration tester")

    client = acme_client.Client(DIRECTORY, key=key, net=net)
    account = client.register(messages.NewRegistration.from_data(email=email))
    client.agree_to_tos(account)
    client.account = account
    return client

class NoClientError(ValueError):
    """
    An error that occurs when no acme.Client is provided to a function that
    requires one.
    """
    pass

class EmailRequiredError(ValueError):
    """
    An error that occurs when a None email is provided to update_email.
    """

def update_email(client, email):
    """
    Use a provided acme.Client to update the client's account to the specified
    email.
    """
    if client is None:
        raise NoClientError("update_email requires a valid acme.Client argument")
    if email is None:
        raise EmailRequiredError("update_email requires an email argument")
    if not email.startswith("mailto:"):
        email = "mailto:"+ email
    acct = client.account
    updatedAcct = acct.update(body=acct.body.update(contact=(email,)))
    return client.update_registration(updatedAcct)

def get_chall(authz, typ):
    for chall_body in authz.body.challenges:
        if isinstance(chall_body.chall, typ):
            return chall_body
    raise Exception("No %s challenge found" % typ)

class ValidationError(Exception):
    """An error that occurs during challenge validation."""
    def __init__(self, domain, problem_type, detail, *args, **kwargs):
        self.domain = domain
        self.problem_type = problem_type
        self.detail = detail

    def __str__(self):
        return "%s: %s: %s" % (self.domain, self.problem_type, self.detail)

def issue(client, authzs, cert_output=None):
    """Given a list of authzs that are being processed by the server,
       wait for them to be ready, then request issuance of a cert with a random
       key for the given domains.

       If cert_output is provided, write the cert as a PEM file to that path."""
    domains = [authz.body.identifier.value for authz in authzs]
    pkey = OpenSSL.crypto.PKey()
    pkey.generate_key(OpenSSL.crypto.TYPE_RSA, 2048)
    csr = OpenSSL.crypto.X509Req()
    csr.add_extensions([
        OpenSSL.crypto.X509Extension(
            'subjectAltName',
            critical=False,
            value=', '.join('DNS:' + d for d in domains).encode()
        ),
    ])
    csr.set_pubkey(pkey)
    csr.set_version(2)
    csr.sign(pkey, 'sha256')

    cert_resource = None
    try:
        cert_resource, _ = client.poll_and_request_issuance(josepy.ComparableX509(csr), authzs)
    except acme_errors.PollError as error:
        # If we get a PollError, pick the first failed authz and turn it into a more
        # useful ValidationError that contains details we can look for in tests.
        for authz in error.updated:
            updated_authz = json.loads(urllib2.urlopen(authz.uri).read())
            domain = authz.body.identifier.value,
            for c in updated_authz['challenges']:
                if 'error' in c:
                    err = c['error']
                    raise ValidationError(domain, err['type'], err['detail'])
        # If none of the authz's had an error, just re-raise.
        raise
    if cert_output is not None:
        pem = OpenSSL.crypto.dump_certificate(OpenSSL.crypto.FILETYPE_PEM,
                                              cert_resource.body)
        with open(cert_output, 'w') as f:
            f.write(pem)
    return cert_resource

def http_01_answer(client, chall_body):
    """Return an HTTP01Resource to server in response to the given challenge."""
    response, validation = chall_body.response_and_validation(client.key)
    return standalone.HTTP01RequestHandler.HTTP01Resource(
          chall=chall_body.chall, response=response,
          validation=validation)

def tls_alpn_01_cert(client, chall_body, domain):
    """Return x509 certificate for tls-alpn-01 challenge"""
    response = chall_body.response(client.key)
    cert, key = response.gen_cert(domain)
    return key, cert

def do_dns_challenges(client, authzs):
    cleanup_hosts = []
    for a in authzs:
        c = get_chall(a, challenges.DNS01)
        name, value = (c.validation_domain_name(a.body.identifier.value),
            c.validation(client.key))
        cleanup_hosts.append(name)
        urllib2.urlopen(SET_TXT,
            data=json.dumps({
                "host": name + ".",
                "value": value,
            })).read()
        client.answer_challenge(c, c.response(client.key))
    def cleanup():
        for host in cleanup_hosts:
            urllib2.urlopen(CLEAR_TXT,
                data=json.dumps({
                    "host": host + ".",
                })).read()
    return cleanup

def do_http_challenges(client, authzs):
    port = 5002
    challs = [get_chall(a, challenges.HTTP01) for a in authzs]
    answers = set([http_01_answer(client, c) for c in challs])
    server = standalone.HTTP01Server(("", port), answers)
    thread = threading.Thread(target=server.serve_forever)
    thread.start()

    # Loop until the HTTP01Server is ready.
    while True:
        try:
            urllib2.urlopen("http://localhost:%d" % port)
            break
        except urllib2.URLError:
            time.sleep(0.1)

    for chall_body in challs:
        client.answer_challenge(chall_body, chall_body.response(client.key))

    def cleanup():
        server.shutdown()
        server.server_close()
        thread.join()
    return cleanup

def do_tlsalpn_challenges(client, authzs):
    port = 5001
    example_key, example_cert = load_example_cert()
    server_certs = {'localhost': (example_key, example_cert)}
    challs = {a.body.identifier.value: get_chall(a, challenges.TLSALPN01)
        for a in authzs}
    chall_certs = {domain: tls_alpn_01_cert(client, c, domain)
        for domain, c in challs.items()}
    # TODO: this won't be needed once acme standalone tls-alpn server serves
    # certs correctly, not only challenge certs.
    chall_certs['localhost'] = (example_key, example_cert)
    server = standalone.TLSALPN01Server(("", port), server_certs, chall_certs)
    thread = threading.Thread(target=server.serve_forever)
    thread.start()

    # Loop until the TLSALPN01Server is ready.
    while True:
        try:
            s = socket.socket()
            s.connect(("localhost", port))
            client_ssl = SSL.Connection(SSL.Context(SSL.TLSv1_METHOD), s)
            client_ssl.set_connect_state()
            client_ssl.set_tlsext_host_name("localhost")
            client_ssl.set_alpn_protos([b'acme-tls/1'])
            client_ssl.do_handshake()
            break
        except (socket.error, SSL.Error):
            time.sleep(0.1)
        finally:
            s.close()

    for chall_body in challs.values():
        client.answer_challenge(chall_body, chall_body.response(client.key))

    def cleanup():
        server.shutdown()
        server.server_close()
        thread.join()
    return cleanup

def load_example_cert():
    keypem = open('test/test-example.key', 'rb').read()
    key = OpenSSL.crypto.load_privatekey(OpenSSL.crypto.FILETYPE_PEM, keypem)
    crtpem = open('test/test-example.pem', 'rb').read()
    cert = OpenSSL.crypto.load_certificate(OpenSSL.crypto.FILETYPE_PEM, crtpem)
    return (key, cert)

def auth_and_issue(domains, chall_type="dns-01", email=None, cert_output=None, client=None):
    """Make authzs for each of the given domains, set up a server to answer the
       challenges in those authzs, tell the ACME server to validate the challenges,
       then poll for the authzs to be ready and issue a cert."""
    if client is None:
        client = make_client(email)
    authzs = [client.request_domain_challenges(d) for d in domains]

    if chall_type == "http-01":
        cleanup = do_http_challenges(client, authzs)
    elif chall_type == "dns-01":
        cleanup = do_dns_challenges(client, authzs)
    elif chall_type == "tls-alpn-01":
        cleanup = do_tlsalpn_challenges(client, authzs)
    else:
        raise Exception("invalid challenge type %s" % chall_type)

    try:
        cert_resource = issue(client, authzs, cert_output)
        client.fetch_chain(cert_resource)
        return cert_resource, authzs
    finally:
        cleanup()

def expect_problem(problem_type, func):
    """Run a function. If it raises a ValidationError or messages.Error that
       contains the given problem_type, return. If it raises no error or the wrong
       error, raise an exception."""
    ok = False
    try:
        func()
    except ValidationError as e:
        if e.problem_type == problem_type:
            ok = True
        else:
            raise
    except messages.Error as e:
        if problem_type in e.__str__():
            ok = True
        else:
            raise
    if not ok:
        raise Exception('Expected %s, got no error' % problem_type)

if __name__ == "__main__":
    domains = sys.argv[1:]
    if len(domains) == 0:
        print __doc__
        sys.exit(0)
    try:
        auth_and_issue(domains)
    except messages.Error, e:
        print e
        sys.exit(1)
