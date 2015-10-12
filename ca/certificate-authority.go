// Copyright 2014 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ca

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/jmhodges/clock"
	"github.com/letsencrypt/boulder/cmd"
	"github.com/letsencrypt/boulder/core"
	blog "github.com/letsencrypt/boulder/log"

	cfsslConfig "github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cloudflare/cfssl/config"
	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cloudflare/cfssl/crypto/pkcs11key"
	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cloudflare/cfssl/helpers"
	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cloudflare/cfssl/ocsp"
	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cloudflare/cfssl/signer"
	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cloudflare/cfssl/signer/local"
)

// This map is used to detect algorithms in crypto/x509 that
// are no longer considered sufficiently strong.
// * No MD2, MD5, or SHA-1
// * No DSA
//
// SHA1WithRSA is allowed because there's still a fair bit of it
// out there, but we should try to remove it soon.
var badSignatureAlgorithms = map[x509.SignatureAlgorithm]bool{
	x509.UnknownSignatureAlgorithm: true,
	x509.MD2WithRSA:                true,
	x509.MD5WithRSA:                true,
	x509.DSAWithSHA1:               true,
	x509.DSAWithSHA256:             true,
	x509.ECDSAWithSHA1:             true,
}

// CertificateAuthorityImpl represents a CA that signs certificates, CRLs, and
// OCSP responses.
type CertificateAuthorityImpl struct {
	profile        string
	Signer         signer.Signer
	OCSPSigner     ocsp.Signer
	SA             core.StorageAuthority
	PA             core.PolicyAuthority
	Publisher      core.Publisher
	Clk            clock.Clock       // TODO(jmhodges): should be private
	Log            *blog.AuditLogger // TODO(bifurcation): should be private
	Prefix         int               // Prepended to the serial number
	ValidityPeriod time.Duration
	NotAfter       time.Time
	MaxNames       int
}

// NewCertificateAuthorityImpl creates a CA that talks to a remote CFSSL
// instance.  (To use a local signer, simply instantiate CertificateAuthorityImpl
// directly.)  Communications with the CA are authenticated with MACs,
// using CFSSL's authenticated signature scheme.  A CA created in this way
// issues for a single profile on the remote signer, which is indicated
// by name in this constructor.
func NewCertificateAuthorityImpl(config cmd.CAConfig, clk clock.Clock, issuerCert string) (*CertificateAuthorityImpl, error) {
	var ca *CertificateAuthorityImpl
	var err error
	logger := blog.GetAuditLogger()
	logger.Notice("Certificate Authority Starting")

	if config.SerialPrefix <= 0 || config.SerialPrefix >= 256 {
		err = errors.New("Must have a positive non-zero serial prefix less than 256 for CA.")
		return nil, err
	}

	// CFSSL requires processing JSON configs through its own LoadConfig, so we
	// serialize and then deserialize.
	cfsslJSON, err := json.Marshal(config.CFSSL)
	if err != nil {
		return nil, err
	}
	cfsslConfigObj, err := cfsslConfig.LoadConfig(cfsslJSON)
	if err != nil {
		return nil, err
	}

	// Load the private key, which can be a file or a PKCS#11 key.
	priv, err := loadKey(config.Key)
	if err != nil {
		return nil, err
	}

	issuer, err := core.LoadCert(issuerCert)
	if err != nil {
		return nil, err
	}

	signer, err := local.NewSigner(priv, issuer, x509.SHA256WithRSA, cfsslConfigObj.Signing)
	if err != nil {
		return nil, err
	}

	if config.LifespanOCSP == "" {
		return nil, errors.New("Config must specify an OCSP lifespan period.")
	}
	lifespanOCSP, err := time.ParseDuration(config.LifespanOCSP)
	if err != nil {
		return nil, err
	}

	// Set up our OCSP signer. Note this calls for both the issuer cert and the
	// OCSP signing cert, which are the same in our case.
	ocspSigner, err := ocsp.NewSigner(issuer, issuer, priv, lifespanOCSP)
	if err != nil {
		return nil, err
	}

	ca = &CertificateAuthorityImpl{
		Signer:     signer,
		OCSPSigner: ocspSigner,
		profile:    config.Profile,
		Prefix:     config.SerialPrefix,
		Clk:        clk,
		Log:        logger,
		NotAfter:   issuer.NotAfter,
	}

	if config.Expiry == "" {
		return nil, errors.New("Config must specify an expiry period.")
	}
	ca.ValidityPeriod, err = time.ParseDuration(config.Expiry)
	if err != nil {
		return nil, err
	}

	ca.MaxNames = config.MaxNames

	return ca, nil
}

func loadKey(keyConfig cmd.KeyConfig) (priv crypto.Signer, err error) {
	if keyConfig.File != "" {
		var keyBytes []byte
		keyBytes, err = ioutil.ReadFile(keyConfig.File)
		if err != nil {
			return nil, fmt.Errorf("Could not read key file %s", keyConfig.File)
		}

		priv, err = helpers.ParsePrivateKeyPEM(keyBytes)
		return
	}

	pkcs11Config := keyConfig.PKCS11
	if pkcs11Config.Module == "" ||
		pkcs11Config.TokenLabel == "" ||
		pkcs11Config.PIN == "" ||
		pkcs11Config.PrivateKeyLabel == "" {
		err = fmt.Errorf("Missing a field in pkcs11Config %#v", pkcs11Config)
		return
	}
	priv, err = pkcs11key.New(pkcs11Config.Module,
		pkcs11Config.TokenLabel, pkcs11Config.PIN, pkcs11Config.PrivateKeyLabel)
	return
}

// GenerateOCSP produces a new OCSP response and returns it
func (ca *CertificateAuthorityImpl) GenerateOCSP(xferObj core.OCSPSigningRequest) ([]byte, error) {
	cert, err := x509.ParseCertificate(xferObj.CertDER)
	if err != nil {
		// AUDIT[ Error Conditions ] 9cc4d537-8534-4970-8665-4b382abe82f3
		ca.Log.AuditErr(err)
		return nil, err
	}

	signRequest := ocsp.SignRequest{
		Certificate: cert,
		Status:      xferObj.Status,
		Reason:      int(xferObj.Reason),
		RevokedAt:   xferObj.RevokedAt,
	}

	ocspResponse, err := ca.OCSPSigner.Sign(signRequest)
	return ocspResponse, err
}

// RevokeCertificate revokes the trust of the Cert referred to by the provided Serial.
func (ca *CertificateAuthorityImpl) RevokeCertificate(serial string, reasonCode core.RevocationCode) (err error) {
	coreCert, err := ca.SA.GetCertificate(serial)
	if err != nil {
		// AUDIT[ Revocation Requests ] 4e85d791-09c0-4ab3-a837-d3d67e945134
		ca.Log.AuditErr(err)
		return err
	}
	cert, err := x509.ParseCertificate(coreCert.DER)
	if err != nil {
		// AUDIT[ Revocation Requests ] 4e85d791-09c0-4ab3-a837-d3d67e945134
		ca.Log.AuditErr(err)
		return err
	}

	signRequest := ocsp.SignRequest{
		Certificate: cert,
		Status:      string(core.OCSPStatusRevoked),
		Reason:      int(reasonCode),
		RevokedAt:   ca.Clk.Now(),
	}
	ocspResponse, err := ca.OCSPSigner.Sign(signRequest)
	if err != nil {
		// AUDIT[ Revocation Requests ] 4e85d791-09c0-4ab3-a837-d3d67e945134
		ca.Log.AuditErr(err)
		return err
	}
	err = ca.SA.MarkCertificateRevoked(serial, ocspResponse, reasonCode)
	return err
}

func uniqueHostNames(csr *x509.CertificateRequest) (commonName string, hostNames []string) {
	hostNames = make([]string, len(csr.DNSNames))
	copy(hostNames, csr.DNSNames)
	if len(csr.Subject.CommonName) > 0 {
		hostNames = append(hostNames, csr.Subject.CommonName)
	}

	hostNames = core.UniqueLowerNames(hostNames)

	if len(csr.Subject.CommonName) > 0 {
		commonName = strings.ToLower(csr.Subject.CommonName)
	} else if len(hostNames) > 0 {
		commonName = hostNames[0]
	}
	return
}

// NewCertificateRequest vets a CSR, and if it's one that we're willing to
// issue under, creates a CertificateRequest so that it can be used for
// issuance later.
func (ca *CertificateAuthorityImpl) NewCertificateRequest(req core.CertificateRequest) (core.CertificateRequest, error) {
	emptyCertRequest := core.CertificateRequest{}

	// Check that enough of the request has been filled out
	if !req.ReadyForCA() {
		err := core.MalformedRequestError("Incomplete certificate request forwarded to CA.")
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return emptyCertRequest, err
	}

	// Parse the CSR
	csr, err := x509.ParseCertificateRequest(req.CSR)
	if err != nil {
		err = core.MalformedRequestError("Malformed CSR.")
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return emptyCertRequest, err
	}

	// Verify that the public key in the CSR is acceptable
	key, ok := csr.PublicKey.(crypto.PublicKey)
	if !ok {
		err = core.MalformedRequestError("Invalid public key in CSR.")
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return emptyCertRequest, err
	}
	if err = core.GoodKey(key); err != nil {
		err = core.MalformedRequestError(fmt.Sprintf("Invalid public key in CSR: %s", err.Error()))
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return emptyCertRequest, err
	}
	if badSignatureAlgorithms[csr.SignatureAlgorithm] {
		err = core.MalformedRequestError("Invalid signature algorithm in CSR")
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return emptyCertRequest, err
	}

	// Verify that the hostnames in the CSR are ones that we will issue for
	// Pull hostnames from CSR
	// Authorization is checked by the RA
	commonName, hostNames := uniqueHostNames(csr)
	if len(hostNames) == 0 {
		err = core.MalformedRequestError("Cannot issue a certificate without a hostname.")
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return emptyCertRequest, err
	}

	// Check that there aren't too many names
	if ca.MaxNames > 0 && len(hostNames) > ca.MaxNames {
		err = core.MalformedRequestError(fmt.Sprintf("Certificate request has %d > %d names", len(hostNames), ca.MaxNames))
		ca.Log.WarningErr(err)
		return emptyCertRequest, err
	}

	// Verify that names are allowed by policy
	identifier := core.AcmeIdentifier{Type: core.IdentifierDNS, Value: commonName}
	if err = ca.PA.WillingToIssue(identifier, req.RegistrationID); err != nil {
		err = core.MalformedRequestError(fmt.Sprintf("Policy forbids issuing for name %s", commonName))
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return emptyCertRequest, err
	}
	for _, name := range hostNames {
		identifier = core.AcmeIdentifier{Type: core.IdentifierDNS, Value: name}
		if err = ca.PA.WillingToIssue(identifier, req.RegistrationID); err != nil {
			err = core.MalformedRequestError(fmt.Sprintf("Policy forbids issuing for name %s", name))
			// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
			ca.Log.AuditErr(err)
			return emptyCertRequest, err
		}
	}

	// Fill in the creation time on the request
	// The expiry time, regID, and CSR are set by the RA on the way in.
	req.Created = ca.Clk.Now()
	req.Status = core.StatusValid

	fullReq, err := ca.SA.NewCertificateRequest(req)
	if err != nil {
		err = core.InternalServerError(fmt.Sprintf("Failed to store certificate request: %v", err))
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return emptyCertRequest, err
	}
	if !fullReq.Complete() {
		err = core.InternalServerError("Storage authority returned an incomplete certificate request")
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return emptyCertRequest, err
	}

	return fullReq, nil
}

// IssueCertificate fetches a certificate request by ID and passes it to a
// goroutine for asynchronous issuance.  This method attempts to catch more
// minor errors; errors in the core signing routines will be reflected by
// making the CertificateRequest invalid.
//
// This method marks the certificateRequest as "pending".  If the certificate
// is successfully issued, the state will move to "valid".  If there is an
// error in issuance, the state will move to "invalid".  That is, the
// "invalid" state on certificate requests indicates that an issuance error
// has occurred (and nothing more).
//
// Because issuance errors at this level are likely to be sensitive, We make
// no attempt to automatically recover from these errors, but instead will
// refuse to issue from the request until the state is "valid" again.
func (ca *CertificateAuthorityImpl) IssueCertificate(requestID string, logEventID string) error {
	logEventResult := "error"
	logEvent := blog.CertificateRequestEvent{
		ID:          logEventID + ".issueCertificate",
		RequestTime: time.Now(),
	}

	// No matter what, log the request
	defer func() {
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditObject(fmt.Sprintf("Certificate request - %s", logEventResult), logEvent)
	}()

	req, err := ca.SA.GetCertificateRequest(requestID)
	if err != nil {
		err = core.InternalServerError("Unable to retrieve certificate request returned from DB")
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return err
	}

	if !req.Complete() {
		err = core.InternalServerError("Incomplete certificate request returned from DB")
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return err
	}

	if req.Status != core.StatusValid {
		err = core.InternalServerError("Cannot issue; certificate request is not valid")
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return err
	}

	// Extract hostnames from the CSR
	csr, err := x509.ParseCertificateRequest(req.CSR)
	if err != nil {
		err = core.InternalServerError("Malformed CSR in cached certificate request")
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return err
	}
	commonName, hostNames := uniqueHostNames(csr)

	// Convert the CSR to PEM
	csrPEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr.Raw,
	}))

	// Assign the validity dates
	notAfter := ca.Clk.Now().Add(ca.ValidityPeriod)

	if ca.NotAfter.Before(notAfter) {
		err = core.InternalServerError("Cannot issue a certificate that expires after the intermediate certificate.")
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return err
	}

	if req.Expires.Before(ca.Clk.Now()) {
		err = core.InternalServerError("Cannot issue on an expired request.")
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditErr(err)
		return err
	}

	err = ca.SA.UpdateCertificateRequestStatus(requestID, core.StatusPending)
	if err != nil {
		return err
	}

	signRequest := signer.SignRequest{
		Request: csrPEM,
		Profile: ca.profile,
		Hosts:   hostNames,
		Subject: &signer.Subject{
			CN: commonName,
		},
	}

	// TODO(rlb): Run this in a goroutine to make it async
	ca.signCertificate(csr, signRequest, requestID, logEventID)

	logEventResult = "successful"
	return nil
}

func (ca *CertificateAuthorityImpl) signCertificate(csr *x509.CertificateRequest, signRequest signer.SignRequest, requestID string, logEventID string) {
	logEventResult := "error"
	logEvent := blog.CertificateRequestEvent{
		ID:          logEventID + ".signCertificate",
		RequestTime: time.Now(),
	}

	// No matter what, log the request
	defer func() {
		// AUDIT[ Certificate Requests ] 11917fa4-10ef-4e0d-9105-bacbe7836a3c
		ca.Log.AuditObject(fmt.Sprintf("Certificate request - %s", logEventResult), logEvent)
	}()

	// Generate a new random serial number
	// We want 136 bits of random number, plus an 8-bit instance id prefix.
	const randBits = 136
	serialBytes := make([]byte, randBits/8+1)
	serialBytes[0] = byte(ca.Prefix)
	_, err := rand.Read(serialBytes[1:])
	if err != nil {
		err = core.InternalServerError(err.Error())
		// AUDIT[ Error Conditions ] 9cc4d537-8534-4970-8665-4b382abe82f3
		ca.Log.Audit(fmt.Sprintf("Serial randomness failed, err=[%v]", err))
		return
	}
	serialHex := hex.EncodeToString(serialBytes)
	serialBigInt := big.NewInt(0)
	serialBigInt = serialBigInt.SetBytes(serialBytes)
	signRequest.Serial = serialBigInt
	logEvent.SerialNumber = serialHex

	certPEM, err := ca.Signer.Sign(signRequest)
	if err != nil {
		err = core.InternalServerError(err.Error())
		// AUDIT[ Error Conditions ] 9cc4d537-8534-4970-8665-4b382abe82f3
		ca.Log.Audit(fmt.Sprintf("Signer failed, rolling back: serial=[%s] err=[%v]", serialHex, err))
		ca.SA.UpdateCertificateRequestStatus(requestID, core.StatusInvalid)
		return
	}

	if len(certPEM) == 0 {
		err = core.InternalServerError("No certificate returned by server")
		// AUDIT[ Error Conditions ] 9cc4d537-8534-4970-8665-4b382abe82f3
		ca.Log.Audit(fmt.Sprintf("PEM empty from Signer, rolling back: serial=[%s] err=[%v]", serialHex, err))
		ca.SA.UpdateCertificateRequestStatus(requestID, core.StatusInvalid)
		return
	}

	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		err = core.InternalServerError("Invalid certificate value returned")
		// AUDIT[ Error Conditions ] 9cc4d537-8534-4970-8665-4b382abe82f3
		ca.Log.Audit(fmt.Sprintf("PEM decode error, aborting and rolling back issuance: pem=[%s] err=[%v]", certPEM, err))
		ca.SA.UpdateCertificateRequestStatus(requestID, core.StatusInvalid)
		return
	}
	certDER := block.Bytes

	// This is one last check for uncaught errors
	if err != nil {
		err = core.InternalServerError(err.Error())
		// AUDIT[ Error Conditions ] 9cc4d537-8534-4970-8665-4b382abe82f3
		ca.Log.Audit(fmt.Sprintf("Uncaught error, aborting and rolling back issuance: pem=[%s] err=[%v]", certPEM, err))
		ca.SA.UpdateCertificateRequestStatus(requestID, core.StatusInvalid)
		return
	}

	// Store the cert with the certificate authority, if provided
	cert, err := ca.SA.AddCertificate(certDER, requestID)
	if err != nil {
		err = core.InternalServerError(err.Error())
		// AUDIT[ Error Conditions ] 9cc4d537-8534-4970-8665-4b382abe82f3
		ca.Log.Audit(fmt.Sprintf("Failed RPC to store at SA, orphaning certificate: pem=[%s] err=[%v]", certPEM, err))
		ca.SA.UpdateCertificateRequestStatus(requestID, core.StatusInvalid)
		return
	}

	// Perform post-issuance checks
	certObj, err := x509.ParseCertificate(certDER)
	if err != nil {
		ca.Log.Warning(fmt.Sprintf("Post-Issuance OCSP failed parsing Certificate: %s", err))
		return
	}

	err = ca.MatchesCSR(cert, csr)
	if err != nil {
		err = core.InternalServerError(err.Error())
		// AUDIT[ Error Conditions ] 9cc4d537-8534-4970-8665-4b382abe82f3
		ca.Log.Audit(fmt.Sprintf("Failed RPC to store at SA, orphaning certificate: pem=[%s] err=[%v]", certPEM, err))
		ca.SA.UpdateCertificateRequestStatus(requestID, core.StatusInvalid)
		return
	}

	// Here's where we finally declare success
	// Now that the certificate has been committed, we can fill out
	// the rest of the log event and mark it as successful
	logEvent.CommonName = certObj.Subject.CommonName
	logEvent.Names = certObj.DNSNames
	logEvent.NotBefore = certObj.NotBefore
	logEvent.NotAfter = certObj.NotAfter
	logEvent.ResponseTime = cert.Issued
	logEventResult = "successful"
	ca.SA.UpdateCertificateRequestStatus(requestID, core.StatusValid)

	// Attempt to generate the OCSP Response now. If this raises an error, it is
	// logged but is not returned to the caller, as an error at this point does
	// not constitute an issuance failure.
	serial := core.SerialToString(certObj.SerialNumber)
	ocspRequest := ocsp.SignRequest{
		Certificate: certObj,
		Status:      string(core.OCSPStatusGood),
	}

	ocspResponse, err := ca.OCSPSigner.Sign(ocspRequest)
	if err != nil {
		ca.Log.Warning(fmt.Sprintf("Post-Issuance OCSP failed signing: %s", err))
		return
	}

	err = ca.SA.UpdateOCSP(serial, ocspResponse)
	if err != nil {
		ca.Log.Warning(fmt.Sprintf("Post-Issuance OCSP failed storing: %s", err))
		return
	}

	// Submit the certificate to any configured CT logs
	ca.Publisher.SubmitToCT(certObj.Raw)

	return
}

// MatchesCSR tests the contents of a generated certificate to make sure
// that the PublicKey, CommonName, and DNSNames match those provided in
// the CSR that was used to generate the certificate. It also checks the
// following fields for:
//		* notBefore is not more than 24 hours ago
//		* BasicConstraintsValid is true
//		* IsCA is false
//		* ExtKeyUsage only contains ExtKeyUsageServerAuth & ExtKeyUsageClientAuth
//		* Subject only contains CommonName & Names
func (ca *CertificateAuthorityImpl) MatchesCSR(cert core.Certificate, csr *x509.CertificateRequest) (err error) {
	parsedCertificate, err := x509.ParseCertificate([]byte(cert.DER))
	if err != nil {
		return
	}

	// Check issued certificate matches what was expected from the CSR
	hostNames := make([]string, len(csr.DNSNames))
	copy(hostNames, csr.DNSNames)
	if len(csr.Subject.CommonName) > 0 {
		hostNames = append(hostNames, csr.Subject.CommonName)
	}
	hostNames = core.UniqueLowerNames(hostNames)

	if !core.KeyDigestEquals(parsedCertificate.PublicKey, csr.PublicKey) {
		err = core.InternalServerError("Generated certificate public key doesn't match CSR public key")
		return
	}
	if len(csr.Subject.CommonName) > 0 && parsedCertificate.Subject.CommonName != csr.Subject.CommonName {
		err = core.InternalServerError("Generated certificate CommonName doesn't match CSR CommonName")
		return
	}
	// Sort both slices of names before comparison.
	parsedNames := parsedCertificate.DNSNames
	sort.Strings(parsedNames)
	sort.Strings(hostNames)
	if !reflect.DeepEqual(parsedNames, hostNames) {
		err = core.InternalServerError("Generated certificate DNSNames don't match CSR DNSNames")
		return
	}
	if !reflect.DeepEqual(parsedCertificate.IPAddresses, csr.IPAddresses) {
		err = core.InternalServerError("Generated certificate IPAddresses don't match CSR IPAddresses")
		return
	}
	if !reflect.DeepEqual(parsedCertificate.EmailAddresses, csr.EmailAddresses) {
		err = core.InternalServerError("Generated certificate EmailAddresses don't match CSR EmailAddresses")
		return
	}
	if len(parsedCertificate.Subject.Country) > 0 || len(parsedCertificate.Subject.Organization) > 0 ||
		len(parsedCertificate.Subject.OrganizationalUnit) > 0 || len(parsedCertificate.Subject.Locality) > 0 ||
		len(parsedCertificate.Subject.Province) > 0 || len(parsedCertificate.Subject.StreetAddress) > 0 ||
		len(parsedCertificate.Subject.PostalCode) > 0 || len(parsedCertificate.Subject.SerialNumber) > 0 {
		err = core.InternalServerError("Generated certificate Subject contains fields other than CommonName or Names")
		return
	}
	now := ca.Clk.Now()
	if now.Sub(parsedCertificate.NotBefore) > time.Hour*24 {
		err = core.InternalServerError(fmt.Sprintf("Generated certificate is back dated %s", now.Sub(parsedCertificate.NotBefore)))
		return
	}
	if !parsedCertificate.BasicConstraintsValid {
		err = core.InternalServerError("Generated certificate doesn't have basic constraints set")
		return
	}
	if parsedCertificate.IsCA {
		err = core.InternalServerError("Generated certificate can sign other certificates")
		return
	}
	if !reflect.DeepEqual(parsedCertificate.ExtKeyUsage, []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}) {
		err = core.InternalServerError("Generated certificate doesn't have correct key usage extensions")
		return
	}

	return
}
