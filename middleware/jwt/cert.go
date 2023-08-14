package jwt

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

var store = certStore{
	roots: x509.NewCertPool(),
	certs: make(map[string]*x509.Certificate),
}

func init() {
	if path := os.Getenv("JWT_ROOT_CERTS"); path != "" {
		certs, err := parseCertsFromPEM(path)
		if err != nil {
			log.Printf("failed to parse root certificates: %v", err)
			os.Exit(2)
		}

		if len(certs) == 0 {
			log.Printf("no root certificates")
			os.Exit(2)
		}

		for _, cert := range certs {
			if cert.IsCA {
				store.roots.AddCert(cert)
			}
		}
	}
}

type certStore struct {
	roots *x509.CertPool
	certs map[string]*x509.Certificate
	sync.RWMutex
}

func (s *certStore) get(fingerprint string) (*x509.Certificate, bool) {
	s.RLock()
	defer s.RUnlock()

	chain, ok := s.certs[fingerprint]

	return chain, ok
}

func (s *certStore) add(cert *x509.Certificate) {
	s.Lock()
	defer s.Unlock()

	s.certs[fingerprintBase64(cert)] = cert
}

func (s *certStore) remove(fingerprint string) {
	s.Lock()
	defer s.Unlock()

	delete(s.certs, fingerprint)
}

func parseCertsFromPEM(path string) ([]*x509.Certificate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var certs []*x509.Certificate
	for {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			break
		}

		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %v", err)
		}

		certs = append(certs, cert)
	}

	return certs, nil
}

// Parses the certificate chain from the JWT x5c header defined in RFC7515.
// The returned slice starts with the leaf certificate followed by the intermediates.
func parseCertsFromToken(token *jwt.Token) ([]*x509.Certificate, error) {
	x5c, ok := token.Header["x5c"].([]interface{})
	if !ok {
		return nil, errors.New("invalid x5c header")
	}
	if len(x5c) == 0 {
		return nil, errors.New("no certificates")
	}

	certs := make([]*x509.Certificate, len(x5c))
	for i, cert := range x5c {
		der, err := base64.StdEncoding.DecodeString(cert.(string))
		if err != nil {
			return nil, fmt.Errorf("failed to decode certificate %d: %v", i, err)
		}

		c, err := x509.ParseCertificate(der)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate %d: %v", i, err)
		}

		if i == 0 && c.IsCA || i > 0 && !c.IsCA {
			return nil, errors.New("invalid certificate chain")
		}

		certs[i] = c
	}

	return certs, nil
}

func verifyCert(leaf *x509.Certificate, intermediates []*x509.Certificate, roots *x509.CertPool) error {
	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: x509.NewCertPool(),
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	for _, cert := range intermediates {
		if !cert.IsCA {
			return errors.New("invalid intermediate certificate")
		}
		opts.Intermediates.AddCert(cert)
	}

	_, err := leaf.Verify(opts)
	if err != nil {
		return fmt.Errorf("failed to verify and build certificate chain: %v", err)
	}

	if leaf.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
		return errors.New("invalid key usage")
	}

	return err
}

func fingerprintBase64(cert *x509.Certificate) string {
	sum := sha256.Sum256(cert.Raw)
	return base64.StdEncoding.EncodeToString(sum[:])
}
