// Copyright 2014 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ca

import (
	//"bytes"
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	cfsslConfig "github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cloudflare/cfssl/config"
	ocspConfig "github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cloudflare/cfssl/ocsp/config"
	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/jmhodges/clock"
	"github.com/letsencrypt/boulder/cmd"
	"github.com/letsencrypt/boulder/mocks"
	"github.com/letsencrypt/boulder/policy"
	"github.com/letsencrypt/boulder/sa/satest"

	"github.com/letsencrypt/boulder/core"
	"github.com/letsencrypt/boulder/sa"
	"github.com/letsencrypt/boulder/test"
)

var (
	CAkeyPEM  = mustRead("./testdata/ca_key.pem")
	CAcertPEM = mustRead("./testdata/ca_cert.pem")

	// CSR generated by Go:
	// * Random public key
	// * CN = not-example.com
	// * DNSNames = not-example.com, www.not-example.com
	CNandSANCSR = mustRead("./testdata/cn_and_san.der.csr")

	// CSR generated by Go:
	// * Random public key
	// * CN = not-example.com
	// * DNSNames = [none]
	NoSANCSR = mustRead("./testdata/no_san.der.csr")

	// CSR generated by Go:
	// * Random public key
	// * C = US
	// * CN = [none]
	// * DNSNames = not-example.com
	NoCNCSR = mustRead("./testdata/no_cn.der.csr")

	// CSR generated by Go:
	// * Random public key
	// * C = US
	// * CN = [none]
	// * DNSNames = [none]
	NoNameCSR = mustRead("./testdata/no_name.der.csr")

	// CSR generated by Go:
	// * Random public key
	// * CN = [none]
	// * DNSNames = a.example.com, a.example.com
	DupeNameCSR = mustRead("./testdata/dupe_name.der.csr")

	// CSR generated by Go:
	// * Random pulic key
	// * CN = [none]
	// * DNSNames = not-example.com, www.not-example.com, mail.example.com
	TooManyNameCSR = mustRead("./testdata/too_many_names.der.csr")

	// CSR generated by Go:
	// * Random public key -- 512 bits long
	// * CN = (none)
	// * DNSNames = not-example.com, www.not-example.com, mail.not-example.com
	ShortKeyCSR = mustRead("./testdata/short_key.der.csr")

	// CSR generated by Go:
	// * Random public key
	// * CN = (none)
	// * DNSNames = not-example.com, www.not-example.com, mail.not-example.com
	// * Signature algorithm: SHA1WithRSA
	BadAlgorithmCSR = mustRead("./testdata/bad_algorithm.der.csr")
	log             = mocks.UseMockLog()
)

// CFSSL config
const profileName = "ee"
const caKeyFile = "../test/test-ca.key"
const caCertFile = "../test/test-ca.pem"
const minWait = 125 * time.Millisecond
const maxWait = 2 * time.Second

const (
	paDBConnStr = "mysql+tcp://boulder@localhost:3306/boulder_policy_test"
	saDBConnStr = "mysql+tcp://boulder@localhost:3306/boulder_sa_test"
)

func mustRead(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("unable to read %#v: %s", path, err))
	}
	return b
}

type testCtx struct {
	sa       core.StorageAuthority
	caConfig cmd.CAConfig
	reg      core.Registration
	pa       core.PolicyAuthority
	fc       clock.FakeClock
	cleanUp  func()
}

func setup(t *testing.T) *testCtx {
	// Create an SA
	dbMap, err := sa.NewDbMap(saDBConnStr)
	if err != nil {
		t.Fatalf("Failed to create dbMap: %s", err)
	}
	fc := clock.NewFake()
	fc.Add(1 * time.Hour)
	ssa, err := sa.NewSQLStorageAuthority(dbMap, fc)
	if err != nil {
		t.Fatalf("Failed to create SA: %s", err)
	}
	saDBCleanUp := test.ResetTestDatabase(t, dbMap.Db)

	paDbMap, err := sa.NewDbMap(paDBConnStr)
	test.AssertNotError(t, err, "Could not construct dbMap")
	pa, err := policy.NewPolicyAuthorityImpl(paDbMap, false)
	test.AssertNotError(t, err, "Couldn't create PADB")
	paDBCleanUp := test.ResetTestDatabase(t, paDbMap.Db)

	cleanUp := func() {
		saDBCleanUp()
		paDBCleanUp()
	}

	// TODO(jmhodges): use of this pkg here is a bug caused by using a real SA
	reg := satest.CreateWorkingRegistration(t, ssa)

	// Create a CA
	caConfig := cmd.CAConfig{
		Profile:      profileName,
		SerialPrefix: 17,
		Key: cmd.KeyConfig{
			File: caKeyFile,
		},
		Expiry:       "8760h",
		LifespanOCSP: "45m",
		MaxNames:     2,
		CFSSL: cfsslConfig.Config{
			Signing: &cfsslConfig.Signing{
				Profiles: map[string]*cfsslConfig.SigningProfile{
					profileName: &cfsslConfig.SigningProfile{
						Usage:     []string{"server auth", "client auth"},
						CA:        false,
						IssuerURL: []string{"http://not-example.com/issuer-url"},
						OCSP:      "http://not-example.com/ocsp",
						CRL:       "http://not-example.com/crl",

						Policies: []cfsslConfig.CertificatePolicy{
							cfsslConfig.CertificatePolicy{
								ID: cfsslConfig.OID(asn1.ObjectIdentifier{2, 23, 140, 1, 2, 1}),
							},
						},
						ExpiryString: "8760h",
						Backdate:     time.Hour,
						CSRWhitelist: &cfsslConfig.CSRWhitelist{
							PublicKeyAlgorithm: true,
							PublicKey:          true,
							SignatureAlgorithm: true,
						},
					},
				},
				Default: &cfsslConfig.SigningProfile{
					ExpiryString: "8760h",
				},
			},
			OCSP: &ocspConfig.Config{
				CACertFile:        caCertFile,
				ResponderCertFile: caCertFile,
				KeyFile:           caKeyFile,
			},
		},
	}
	return &testCtx{ssa, caConfig, reg, pa, fc, cleanUp}
}

func (ctx *testCtx) importCSR(csr []byte) core.CertificateRequest {
	return core.CertificateRequest{
		RegistrationID: ctx.reg.ID,
		Created:        ctx.fc.Now(),
		Expires:        ctx.fc.Now().AddDate(1, 0, 0),
		Status:         core.StatusValid,
		CSR:            csr,
	}
}

// Use this function to wait for the IssueCertificate->sign->SA handoff
//
//   Note: No need to use this for negative tests ("Test that the CA
//   rejects...").  If IssueCertificate() returns without an error, then
//   the CA has agreed to issue the certificate, which is already a
//   failure for these cases.
func (ctx *testCtx) attemptToIssue(t *testing.T, ca *CertificateAuthorityImpl, csr []byte) (cert core.Certificate, found bool) {
	req := ctx.importCSR(csr)
	req, err := ca.NewCertificateRequest(req)
	test.AssertNotError(t, err, "Failed to import CSR")

	err = ca.IssueCertificate(req.ID, "bogusLogEvent")
	test.AssertNotError(t, err, "CA did not agree to issue under this request")

	found = false
	wait := minWait
	for wait < maxWait {
		time.Sleep(wait)
		wait = 2 * wait

		req, err := ctx.sa.GetCertificateRequest(req.ID)
		test.AssertNotError(t, err, "Unable to retrieve supposedly pending certificate request")
		if err != nil {
			return
		}

		switch req.Status {
		case core.StatusValid:
			cert, err = ctx.sa.GetLatestCertificateForRequest(req.ID)
			test.AssertNotError(t, err, "Error retrieving issued certificate")
			found = (err == nil)
			return
		case core.StatusInvalid:
			test.Assert(t, false, "Issuance failed; see logs for details")
			return
		}
	}
	return
}

func TestFailNoSerial(t *testing.T) {
	ctx := setup(t)
	defer ctx.cleanUp()

	ctx.caConfig.SerialPrefix = 0
	_, err := NewCertificateAuthorityImpl(ctx.caConfig, ctx.fc, caCertFile)
	test.AssertError(t, err, "CA should have failed with no SerialPrefix")
}

func TestRevoke(t *testing.T) {
	ctx := setup(t)
	defer ctx.cleanUp()
	ca, err := NewCertificateAuthorityImpl(ctx.caConfig, ctx.fc, caCertFile)
	test.AssertNotError(t, err, "Failed to create CA")
	ca.PA = ctx.pa
	ca.SA = ctx.sa
	ca.Publisher = &mocks.MockPublisher{}

	certObj, found := ctx.attemptToIssue(t, ca, CNandSANCSR)
	test.Assert(t, found, "Promised cert failed to appear")

	cert, err := x509.ParseCertificate(certObj.DER)
	test.AssertNotError(t, err, "Certificate failed to parse")
	serialString := core.SerialToString(cert.SerialNumber)

	beforeRevoke, err := ctx.sa.GetCertificateStatus(serialString)
	test.AssertNotError(t, err, "Failed to get cert status")

	ctx.fc.Add(1 * time.Hour)

	err = ca.RevokeCertificate(serialString, 0)
	test.AssertNotError(t, err, "Revocation failed")

	status, err := ctx.sa.GetCertificateStatus(serialString)
	test.AssertNotError(t, err, "Failed to get cert status")

	test.AssertEquals(t, status.Status, core.OCSPStatusRevoked)

	if !ctx.fc.Now().Equal(status.OCSPLastUpdated) {
		t.Errorf("OCSPLastUpdated, expected %s, got %s",
			ctx.fc.Now(),
			status.OCSPLastUpdated)
	}
	if !status.OCSPLastUpdated.After(beforeRevoke.OCSPLastUpdated) {
		t.Errorf("OCSPLastUpdated, before revocation: %s; after: %s", beforeRevoke.OCSPLastUpdated, status.OCSPLastUpdated)
	}

}

func TestIssueCertificate(t *testing.T) {
	ctx := setup(t)
	defer ctx.cleanUp()
	ca, err := NewCertificateAuthorityImpl(ctx.caConfig, ctx.fc, caCertFile)
	test.AssertNotError(t, err, "Failed to create CA")
	ca.Publisher = &mocks.MockPublisher{}
	ca.PA = ctx.pa
	ca.SA = ctx.sa

	/*
		  // Uncomment to test with a local signer
			signer, _ := local.NewSigner(caKey, caCert, x509.SHA256WithRSA, nil)
			ca := CertificateAuthorityImpl{
				Signer: signer,
				SA:     sa,
			}
	*/

	csrs := [][]byte{CNandSANCSR, NoSANCSR, NoCNCSR}
	for _, csrDER := range csrs {
		// Sign CSR
		issuedCert, found := ctx.attemptToIssue(t, ca, csrDER)
		test.Assert(t, found, "Promised cert failed to appear")

		// Verify cert contents
		cert, err := x509.ParseCertificate(issuedCert.DER)
		test.AssertNotError(t, err, "Certificate failed to parse")

		test.AssertEquals(t, cert.Subject.CommonName, "not-example.com")

		switch len(cert.DNSNames) {
		case 1:
			if cert.DNSNames[0] != "not-example.com" {
				t.Errorf("Improper list of domain names %v", cert.DNSNames)
			}
		case 2:
			switch {
			case (cert.DNSNames[0] == "not-example.com" && cert.DNSNames[1] == "www.not-example.com"):
				t.Log("case 1")
			case (cert.DNSNames[0] == "www.not-example.com" && cert.DNSNames[1] == "not-example.com"):
				t.Log("case 2")
			default:
				t.Errorf("Improper list of domain names %v", cert.DNSNames)
			}

		default:
			t.Errorf("Improper list of domain names %v", cert.DNSNames)
		}

		// Test is broken by CFSSL Issue #156
		// https://github.com/cloudflare/cfssl/issues/156
		if len(cert.Subject.Country) > 0 {
			// Uncomment the Errorf as soon as upstream #156 is fixed
			// t.Errorf("Subject contained unauthorized values: %v", cert.Subject)
			t.Logf("Subject contained unauthorized values: %v", cert.Subject)
		}

		// Verify that the cert got stored in the DB
		serialString := core.SerialToString(cert.SerialNumber)
		storedCert, err := ctx.sa.GetCertificate(serialString)
		test.AssertNotError(t, err,
			fmt.Sprintf("Certificate %s not found in database", serialString))
		test.AssertByteEquals(t, issuedCert.DER, storedCert.DER)

		certStatus, err := ctx.sa.GetCertificateStatus(serialString)
		test.AssertNotError(t, err,
			fmt.Sprintf("Error fetching status for certificate %s", serialString))
		test.Assert(t, certStatus.Status == core.OCSPStatusGood, "Certificate status was not good")
		test.Assert(t, certStatus.SubscriberApproved == false, "Subscriber shouldn't have approved cert yet.")
	}
}

func TestRejectNoName(t *testing.T) {
	ctx := setup(t)
	defer ctx.cleanUp()
	ca, err := NewCertificateAuthorityImpl(ctx.caConfig, ctx.fc, caCertFile)
	test.AssertNotError(t, err, "Failed to create CA")
	ca.Publisher = &mocks.MockPublisher{}
	ca.PA = ctx.pa
	ca.SA = ctx.sa

	// Test that the CA rejects CSRs with no names
	req := ctx.importCSR(NoNameCSR)
	_, err = ca.NewCertificateRequest(req)
	test.AssertError(t, err, "CA improperly agreed to create a certificate with no name")
	_, ok := err.(core.MalformedRequestError)
	test.Assert(t, ok, "Incorrect error type returned")
}

func TestRejectTooManyNames(t *testing.T) {
	ctx := setup(t)
	defer ctx.cleanUp()
	ca, err := NewCertificateAuthorityImpl(ctx.caConfig, ctx.fc, caCertFile)
	test.AssertNotError(t, err, "Failed to create CA")
	ca.Publisher = &mocks.MockPublisher{}
	ca.PA = ctx.pa
	ca.SA = ctx.sa

	// Test that the CA rejects a CSR with too many names
	req := ctx.importCSR(TooManyNameCSR)
	_, err = ca.NewCertificateRequest(req)
	test.AssertError(t, err, "Issued certificate with too many names")
	_, ok := err.(core.MalformedRequestError)
	test.Assert(t, ok, "Incorrect error type returned")
}

func TestDeduplication(t *testing.T) {
	ctx := setup(t)
	defer ctx.cleanUp()
	ca, err := NewCertificateAuthorityImpl(ctx.caConfig, ctx.fc, caCertFile)
	test.AssertNotError(t, err, "Failed to create CA")
	ca.Publisher = &mocks.MockPublisher{}
	ca.PA = ctx.pa
	ca.SA = ctx.sa

	// Test that the CA collapses duplicate names
	cert, found := ctx.attemptToIssue(t, ca, DupeNameCSR)
	test.Assert(t, found, "Promised cert failed to appear")

	parsedCert, err := x509.ParseCertificate(cert.DER)
	test.AssertNotError(t, err, "Error parsing certificate produced by CA")

	correctName := "a.not-example.com"
	correctNames := len(parsedCert.DNSNames) == 1 &&
		parsedCert.DNSNames[0] == correctName &&
		parsedCert.Subject.CommonName == correctName
	fmt.Println("Names in the cert:", parsedCert.DNSNames)
	test.Assert(t, correctNames, "Incorrect set of names in deduplicated certificate")
}

func TestRejectValidityTooLong(t *testing.T) {
	ctx := setup(t)
	defer ctx.cleanUp()
	ca, err := NewCertificateAuthorityImpl(ctx.caConfig, ctx.fc, caCertFile)
	test.AssertNotError(t, err, "Failed to create CA")
	ca.Publisher = &mocks.MockPublisher{}
	ca.PA = ctx.pa
	ca.SA = ctx.sa

	// Test that the CA rejects CSRs that would expire after the intermediate cert
	ca.NotAfter = ctx.fc.Now()
	req := ctx.importCSR(NoCNCSR)
	req, err = ca.NewCertificateRequest(req)
	test.AssertNotError(t, err, "Failed to import CSR")

	err = ca.IssueCertificate(req.ID, "bogusLogEvent")
	test.AssertEquals(t, err.Error(), "Cannot issue a certificate that expires after the intermediate certificate.")
	_, ok := err.(core.InternalServerError)
	test.Assert(t, ok, "Incorrect error type returned")
}

func TestShortKey(t *testing.T) {
	ctx := setup(t)
	defer ctx.cleanUp()
	ca, err := NewCertificateAuthorityImpl(ctx.caConfig, ctx.fc, caCertFile)
	ca.Publisher = &mocks.MockPublisher{}
	ca.PA = ctx.pa
	ca.SA = ctx.sa

	// Test that the CA rejects CSRs that would expire after the intermediate cert
	req := ctx.importCSR(ShortKeyCSR)
	_, err = ca.NewCertificateRequest(req)
	test.AssertError(t, err, "Issued a certificate with too short a key.")
	_, ok := err.(core.MalformedRequestError)
	test.Assert(t, ok, "Incorrect error type returned")
}

func TestRejectBadAlgorithm(t *testing.T) {
	ctx := setup(t)
	defer ctx.cleanUp()
	ca, err := NewCertificateAuthorityImpl(ctx.caConfig, ctx.fc, caCertFile)
	ca.Publisher = &mocks.MockPublisher{}
	ca.PA = ctx.pa
	ca.SA = ctx.sa

	// Test that the CA rejects CSRs that would expire after the intermediate cert
	req := ctx.importCSR(BadAlgorithmCSR)
	_, err = ca.NewCertificateRequest(req)
	test.AssertError(t, err, "Issued a certificate based on a CSR with a weak algorithm.")
	_, ok := err.(core.MalformedRequestError)
	test.Assert(t, ok, "Incorrect error type returned")
}
