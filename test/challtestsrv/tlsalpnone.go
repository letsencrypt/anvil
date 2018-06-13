package challtestsrv

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"net/http"
	"time"

	"github.com/letsencrypt/boulder/va"
)

// AddTLSALPNChallenge adds a new TLS-ALPN-01 key authorization for the given host
func (s *ChallSrv) AddTLSALPNChallenge(host, content string) {
	s.challMu.Lock()
	defer s.challMu.Unlock()
	s.taOne[host] = content
}

// DeleteTLSALPNChallenge deletes the key authorization for a given host
func (s *ChallSrv) DeleteTLSALPNChallenge(host string) {
	s.challMu.Lock()
	defer s.challMu.Unlock()
	if _, ok := s.taOne[host]; ok {
		delete(s.taOne, host)
	}
}

// GetTLSALPNChallenge returns the TLS-ALPN-01 key authorization for the given host
// (if it exists) and a true bool. If the host does not exist then an empty
// string and a false bool are returned.
func (s *ChallSrv) GetTLSALPNChallenge(host string) (string, bool) {
	s.challMu.RLock()
	defer s.challMu.RUnlock()
	content, present := s.taOne[host]
	return content, present
}

func (s *ChallSrv) ServeChallengeCertFunc(k *ecdsa.PrivateKey) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if len(hello.SupportedProtos) != 0 || hello.SupportedProtos[0] != va.ACMETLS1Protocol {
			return nil, fmt.Errorf("ALPN failed, ClientHelloInfo.SupportedProtos: %s", hello.SupportedProtos)
		}

		ka, found := s.GetTLSALPNChallenge(hello.ServerName)
		if !found {
			return nil, fmt.Errorf("unknown ClientHelloInfo.ServerName: %s", hello.ServerName)
		}

		kaHash := sha256.Sum256([]byte(ka))
		extValue, err := asn1.Marshal(kaHash[:])
		if err != nil {
			return nil, fmt.Errorf("failed marshalling hash OCTET STRING: %s", err)
		}
		certTmpl := x509.Certificate{
			DNSNames: []string{hello.ServerName},
			Extensions: []pkix.Extension{
				{
					Id:       va.IdPeAcmeIdentifierV1,
					Critical: true,
					Value:    extValue,
				},
			},
		}
		certBytes, err := x509.CreateCertificate(rand.Reader, &certTmpl, &certTmpl, k.Public(), k)
		if err != nil {
			return nil, fmt.Errorf("failed creating challenge certificate: %s", err)
		}
		return &tls.Certificate{
			Certificate: [][]byte{certBytes},
			PrivateKey:  k,
		}, nil
	}
}

type challTLSServer struct {
	*http.Server
}

func (c challTLSServer) Shutdown() error {
	return c.Server.Shutdown(context.Background())
}

func (c challTLSServer) ListenAndServe() error {
	// We never want to serve a plain cert so leave certFile and keyFile
	// empty. If we don't know the SNI name/ALPN fails the handshake will
	// fail anyway.
	return c.Server.ListenAndServeTLS("", "")
}

func taOneServer(address string, challSrv *ChallSrv) challengeServer {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	srv := &http.Server{
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		TLSConfig: &tls.Config{
			NextProtos:     []string{va.ACMETLS1Protocol},
			GetCertificate: challSrv.ServeChallengeCertFunc(key),
		},
	}
	srv.SetKeepAlivesEnabled(false)
	return challTLSServer{srv}
}
