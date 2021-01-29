package publisher

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	ct "github.com/google/certificate-transparency-go"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/letsencrypt/boulder/core"
	"github.com/letsencrypt/boulder/issuance"
	blog "github.com/letsencrypt/boulder/log"
	"github.com/letsencrypt/boulder/metrics"
	pubpb "github.com/letsencrypt/boulder/publisher/proto"
	"github.com/letsencrypt/boulder/test"
)

var log = blog.UseMock()
var ctx = context.Background()

func getPort(srvURL string) (int, error) {
	url, err := url.Parse(srvURL)
	if err != nil {
		return 0, err
	}
	_, portString, err := net.SplitHostPort(url.Host)
	if err != nil {
		return 0, err
	}
	port, err := strconv.ParseInt(portString, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(port), nil
}

type testLogSrv struct {
	*httptest.Server
	submissions int64
}

func logSrv(k *ecdsa.PrivateKey) *testLogSrv {
	testLog := &testLogSrv{}
	m := http.NewServeMux()
	m.HandleFunc("/ct/", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var jsonReq ctSubmissionRequest
		err := decoder.Decode(&jsonReq)
		if err != nil {
			return
		}
		precert := false
		if r.URL.Path == "/ct/v1/add-pre-chain" {
			precert = true
		}
		sct := CreateTestingSignedSCT(jsonReq.Chain, k, precert, time.Now())
		fmt.Fprint(w, string(sct))
		atomic.AddInt64(&testLog.submissions, 1)
	})

	testLog.Server = httptest.NewUnstartedServer(m)
	testLog.Server.Start()
	return testLog
}

// lyingLogSrv always signs SCTs with the timestamp it was given.
func lyingLogSrv(k *ecdsa.PrivateKey, timestamp time.Time) *testLogSrv {
	testLog := &testLogSrv{}
	m := http.NewServeMux()
	m.HandleFunc("/ct/", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var jsonReq ctSubmissionRequest
		err := decoder.Decode(&jsonReq)
		if err != nil {
			return
		}
		precert := false
		if r.URL.Path == "/ct/v1/add-pre-chain" {
			precert = true
		}
		sct := CreateTestingSignedSCT(jsonReq.Chain, k, precert, timestamp)
		fmt.Fprint(w, string(sct))
		atomic.AddInt64(&testLog.submissions, 1)
	})

	testLog.Server = httptest.NewUnstartedServer(m)
	testLog.Server.Start()
	return testLog
}

func errorBodyLogSrv() *httptest.Server {
	m := http.NewServeMux()
	m.HandleFunc("/ct/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("well this isn't good now is it."))
	})

	server := httptest.NewUnstartedServer(m)
	server.Start()
	return server
}

func setup(t *testing.T) (*Impl, []*x509.Certificate, *ecdsa.PrivateKey) {
	// Load a chain with two intermediates
	testCA2Cert, err := core.LoadCert("../test/test-ca2.pem")
	test.AssertNotError(t, err, "Unable to load ../test/test-ca2.pem")

	testCA2CrossCert, err := core.LoadCert("../test/test-ca2-cross.pem")
	test.AssertNotError(t, err, "Unable to load ../test/test-ca2-cross.pem")

	// Create our ct.ASN1Cert bundles mapped to their corresponding
	// IssueNameID
	issuerBundles := map[issuance.IssuerNameID][]ct.ASN1Cert{
		// Chain with two intermediates
		issuance.IssuerNameID(18337263084599622): {
			ct.ASN1Cert{Data: testCA2Cert.Raw},
			ct.ASN1Cert{Data: testCA2CrossCert.Raw},
		},
		// Chain with two intermediates mapped to the IssuerNameID of
		// publisher/test/178.pem
		issuance.IssuerNameID(66191037641995744): {
			ct.ASN1Cert{Data: testCA2Cert.Raw},
			ct.ASN1Cert{Data: testCA2CrossCert.Raw},
		},
	}
	pub := New(
		issuerBundles,
		"test-user-agent/1.0",
		log,
		metrics.NoopRegisterer)

	// Load leaf that corresponds to our chain with two intermediates
	testEECert, err := core.LoadCert("../test/test-ee.pem")
	test.AssertNotError(t, err, "Unable to load ../test/test-ee.pem")

	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	test.AssertNotError(t, err, "Couldn't generate test key")

	return pub, []*x509.Certificate{testEECert}, k
}

func addLog(t *testing.T, pub *Impl, port int, pubKey *ecdsa.PublicKey) *Log {
	uri := fmt.Sprintf("http://localhost:%d", port)
	der, err := x509.MarshalPKIXPublicKey(pubKey)
	test.AssertNotError(t, err, "Failed to marshal key")
	newLog, err := NewLog(uri, base64.StdEncoding.EncodeToString(der), "test-user-agent/1.0", log)
	test.AssertNotError(t, err, "Couldn't create log")
	test.AssertEquals(t, newLog.uri, fmt.Sprintf("http://localhost:%d", port))
	return newLog
}

func makePrecert(k *ecdsa.PrivateKey) (map[issuance.IssuerNameID][]ct.ASN1Cert, []byte, error) {
	rootTmpl := x509.Certificate{
		SerialNumber:          big.NewInt(0),
		Subject:               pkix.Name{CommonName: "root"},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	rootBytes, err := x509.CreateCertificate(rand.Reader, &rootTmpl, &rootTmpl, k.Public(), k)
	if err != nil {
		return nil, nil, err
	}
	root, err := x509.ParseCertificate(rootBytes)
	if err != nil {
		return nil, nil, err
	}
	precertTmpl := x509.Certificate{
		SerialNumber: big.NewInt(0),
		ExtraExtensions: []pkix.Extension{
			{Id: asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 11129, 2, 4, 3}, Critical: true, Value: []byte{0x05, 0x00}},
		},
	}
	precert, err := x509.CreateCertificate(rand.Reader, &precertTmpl, root, k.Public(), k)
	if err != nil {
		return nil, nil, err
	}
	precertX509, err := x509.ParseCertificate(precert)
	if err != nil {
		return nil, nil, err
	}
	precertIssuerNameID := issuance.GetIssuerNameID(precertX509)
	bundles := map[issuance.IssuerNameID][]ct.ASN1Cert{
		precertIssuerNameID: {
			ct.ASN1Cert{Data: rootBytes},
		},
	}
	return bundles, precert, err
}

func TestTimestampVerificationFuture(t *testing.T) {
	pub, _, k := setup(t)

	server := lyingLogSrv(k, time.Now().Add(time.Hour))
	defer server.Close()
	port, err := getPort(server.URL)
	test.AssertNotError(t, err, "Failed to get test server port")
	testLog := addLog(t, pub, port, &k.PublicKey)

	// Precert
	issuerBundles, precert, err := makePrecert(k)
	test.AssertNotError(t, err, "Failed to create test leaf")
	pub.issuerBundles = issuerBundles

	_, err = pub.SubmitToSingleCTWithResult(ctx, &pubpb.Request{LogURL: testLog.uri, LogPublicKey: testLog.logID, Der: precert, Precert: true})
	if err == nil {
		t.Fatal("Expected error for lying log server, got none")
	}
	if !strings.HasPrefix(err.Error(), "SCT Timestamp was too far in the future") {
		t.Fatalf("Got wrong error: %s", err)
	}
}

func TestTimestampVerificationPast(t *testing.T) {
	pub, _, k := setup(t)

	server := lyingLogSrv(k, time.Now().Add(-time.Hour))
	defer server.Close()
	port, err := getPort(server.URL)
	test.AssertNotError(t, err, "Failed to get test server port")
	testLog := addLog(t, pub, port, &k.PublicKey)

	// Precert
	issuerBundles, precert, err := makePrecert(k)
	test.AssertNotError(t, err, "Failed to create test leaf")

	pub.issuerBundles = issuerBundles

	_, err = pub.SubmitToSingleCTWithResult(ctx, &pubpb.Request{LogURL: testLog.uri, LogPublicKey: testLog.logID, Der: precert, Precert: true})
	if err == nil {
		t.Fatal("Expected error for lying log server, got none")
	}
	if !strings.HasPrefix(err.Error(), "SCT Timestamp was too far in the past") {
		t.Fatalf("Got wrong error: %s", err)
	}
}

func TestLogCache(t *testing.T) {
	cache := logCache{
		logs: make(map[string]*Log),
	}

	// Adding a log with an invalid base64 public key should error
	_, err := cache.AddLog("www.test.com", "1234", "test-user-agent/1.0", log)
	test.AssertError(t, err, "AddLog() with invalid base64 pk didn't error")

	// Adding a log with an invalid URI should error
	_, err = cache.AddLog(":", "", "test-user-agent/1.0", log)
	test.AssertError(t, err, "AddLog() with an invalid log URI didn't error")

	// Create one keypair & base 64 public key
	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	test.AssertNotError(t, err, "ecdsa.GenerateKey() failed for k1")
	der1, err := x509.MarshalPKIXPublicKey(&k1.PublicKey)
	test.AssertNotError(t, err, "x509.MarshalPKIXPublicKey(der1) failed")
	k1b64 := base64.StdEncoding.EncodeToString(der1)

	// Create a second keypair & base64 public key
	k2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	test.AssertNotError(t, err, "ecdsa.GenerateKey() failed for k2")
	der2, err := x509.MarshalPKIXPublicKey(&k2.PublicKey)
	test.AssertNotError(t, err, "x509.MarshalPKIXPublicKey(der2) failed")
	k2b64 := base64.StdEncoding.EncodeToString(der2)

	// Adding the first log should not produce an error
	l1, err := cache.AddLog("http://log.one.example.com", k1b64, "test-user-agent/1.0", log)
	test.AssertNotError(t, err, "cache.AddLog() failed for log 1")
	test.AssertEquals(t, cache.Len(), 1)
	test.AssertEquals(t, l1.uri, "http://log.one.example.com")
	test.AssertEquals(t, l1.logID, k1b64)

	// Adding it again should not produce any errors, or increase the Len()
	l1, err = cache.AddLog("http://log.one.example.com", k1b64, "test-user-agent/1.0", log)
	test.AssertNotError(t, err, "cache.AddLog() failed for second add of log 1")
	test.AssertEquals(t, cache.Len(), 1)
	test.AssertEquals(t, l1.uri, "http://log.one.example.com")
	test.AssertEquals(t, l1.logID, k1b64)

	// Adding a second log should not error and should increase the Len()
	l2, err := cache.AddLog("http://log.two.example.com", k2b64, "test-user-agent/1.0", log)
	test.AssertNotError(t, err, "cache.AddLog() failed for log 2")
	test.AssertEquals(t, cache.Len(), 2)
	test.AssertEquals(t, l2.uri, "http://log.two.example.com")
	test.AssertEquals(t, l2.logID, k2b64)
}

func TestLogErrorBody(t *testing.T) {
	pub, leaves, k := setup(t)

	srv := errorBodyLogSrv()
	defer srv.Close()
	port, err := getPort(srv.URL)
	test.AssertNotError(t, err, "Failed to get test server port")

	log.Clear()
	logURI := fmt.Sprintf("http://localhost:%d", port)
	pkDER, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	test.AssertNotError(t, err, "Failed to marshal key")
	pkB64 := base64.StdEncoding.EncodeToString(pkDER)
	_, err = pub.SubmitToSingleCTWithResult(context.Background(), &pubpb.Request{
		LogURL:       logURI,
		LogPublicKey: pkB64,
		Der:          leaves[0].Raw,
	})
	test.AssertError(t, err, "SubmitToSingleCTWithResult didn't fail")
	test.AssertEquals(t, len(log.GetAllMatching("well this isn't good now is it")), 1)
}

func TestHTTPStatusMetric(t *testing.T) {
	pub, leaves, k := setup(t)

	badSrv := errorBodyLogSrv()
	defer badSrv.Close()
	port, err := getPort(badSrv.URL)
	test.AssertNotError(t, err, "Failed to get test server port")
	logURI := fmt.Sprintf("http://localhost:%d", port)

	pkDER, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	test.AssertNotError(t, err, "Failed to marshal key")
	pkB64 := base64.StdEncoding.EncodeToString(pkDER)
	_, err = pub.SubmitToSingleCTWithResult(context.Background(), &pubpb.Request{
		LogURL:       logURI,
		LogPublicKey: pkB64,
		Der:          leaves[0].Raw,
	})
	test.AssertError(t, err, "SubmitToSingleCTWithResult didn't fail")
	test.AssertEquals(t, test.CountHistogramSamples(pub.metrics.submissionLatency.With(prometheus.Labels{
		"log":         logURI,
		"status":      "error",
		"http_status": "400",
	})), 1)

	pub, leaves, k = setup(t)
	pkDER, err = x509.MarshalPKIXPublicKey(&k.PublicKey)
	test.AssertNotError(t, err, "Failed to marshal key")
	pkB64 = base64.StdEncoding.EncodeToString(pkDER)
	workingSrv := logSrv(k)
	defer workingSrv.Close()
	port, err = getPort(workingSrv.URL)
	test.AssertNotError(t, err, "Failed to get test server port")
	logURI = fmt.Sprintf("http://localhost:%d", port)

	_, err = pub.SubmitToSingleCTWithResult(context.Background(), &pubpb.Request{
		LogURL:       logURI,
		LogPublicKey: pkB64,
		Der:          leaves[0].Raw,
	})
	test.AssertNotError(t, err, "SubmitToSingleCTWithResult failed")
	test.AssertEquals(t, test.CountHistogramSamples(pub.metrics.submissionLatency.With(prometheus.Labels{
		"log":         logURI,
		"status":      "success",
		"http_status": "",
	})), 1)
}
