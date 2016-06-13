package creds

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/letsencrypt/boulder/test"
)

func slowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Millisecond * 200)
}

func TestTransportCredentials(t *testing.T) {
	serverA := httptest.NewTLSServer(nil)
	defer serverA.Close()
	addrA := serverA.Listener.Addr().String()
	serverB := httptest.NewTLSServer(nil)
	defer serverB.Close()
	addrB := serverB.Listener.Addr().String()

	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	test.AssertNotError(t, err, "rsa.GenerateKey failed")

	temp := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "A",
		},
		NotBefore:             time.Unix(1000, 0),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		BasicConstraintsValid: true,
		IsCA: true,
	}
	derA, err := x509.CreateCertificate(rand.Reader, temp, temp, priv.Public(), priv)
	test.AssertNotError(t, err, "x509.CreateCertificate failed")
	certA, err := x509.ParseCertificate(derA)
	test.AssertNotError(t, err, "x509.ParserCertificate failed")
	temp.Subject.CommonName = "B"
	derB, err := x509.CreateCertificate(rand.Reader, temp, temp, priv.Public(), priv)
	test.AssertNotError(t, err, "x509.CreateCertificate failed")
	certB, err := x509.ParseCertificate(derB)
	test.AssertNotError(t, err, "x509.ParserCertificate failed")
	// XXX: don't do this
	serverA.TLS.Certificates = []tls.Certificate{{Certificate: [][]byte{derA}, PrivateKey: priv}}
	serverB.TLS.Certificates = []tls.Certificate{{Certificate: [][]byte{derB}, PrivateKey: priv}}

	roots := x509.NewCertPool()
	roots.AddCert(certA)
	roots.AddCert(certB)

	tc, err := New([]string{"A:2020", "B:3030"}, roots, nil)
	test.AssertNotError(t, err, "New failed")

	rawConnA, err := net.Dial("tcp", addrA)
	test.AssertNotError(t, err, "net.Dial failed")
	defer func() {
		_ = rawConnA.Close()
	}()

	conn, _, err := tc.ClientHandshake("A:2020", rawConnA, 0)
	test.AssertNotError(t, err, "tc.ClientHandshake failed")
	test.Assert(t, conn != nil, "tc.ClientHandshake returned a nil net.Conn")

	rawConnB, err := net.Dial("tcp", addrB)
	test.AssertNotError(t, err, "net.Dial failed")
	defer func() {
		_ = rawConnB.Close()
	}()

	conn, _, err = tc.ClientHandshake("B:3030", rawConnB, 0)
	test.AssertNotError(t, err, "tc.ClientHandshake failed")
	test.Assert(t, conn != nil, "tc.ClientHandshake returned a nil net.Conn")

	ln, err := net.Listen("tcp", ":")
	test.AssertNotError(t, err, "net.Listen failed")
	addrC := ln.Addr().String()
	go func() {
		for {
			_, err := ln.Accept()
			test.AssertNotError(t, err, "ln.Accept failed")
			time.Sleep(time.Second)
		}
	}()

	rawConnC, err := net.Dial("tcp", addrC)
	test.AssertNotError(t, err, "net.Dial failed")
	defer func() {
		_ = rawConnB.Close()
	}()

	conn, _, err = tc.ClientHandshake("A:2020", rawConnC, time.Millisecond)
	test.AssertError(t, err, "tc.ClientHandshake didn't timeout")
	test.AssertEquals(t, err.Error(), "boulder/grpc/creds: TLS handshake timed out")
	test.Assert(t, conn == nil, "tc.ClientHandshake returned a non-nil net.Conn on failure")
}
