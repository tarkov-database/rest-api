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
	"slices"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

var store certStore = certStore{
	roots: x509.NewCertPool(),
	certs: make(map[string]*certChain),
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

	if path := os.Getenv("JWT_CERTS"); path != "" {
		certs, err := parseCertsFromPEM(path)
		if err != nil {
			log.Printf("failed to parse certificates: %v", err)
			os.Exit(2)
		}

		if len(certs) == 0 {
			log.Printf("no certificates")
			os.Exit(2)
		}

		lastIdx := len(certs) - 1

		intermediates := x509.NewCertPool()
		for _, cert := range certs[:lastIdx] {
			intermediates.AddCert(cert)
		}

		store.add(&certChain{
			intermediates: intermediates,
			leaf:          certs[lastIdx],
		})
	}
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

		if block.Type != "CERTIFICATE" {
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

type certStore struct {
	roots *x509.CertPool
	certs map[string]*certChain
	sync.RWMutex
}

func (s *certStore) add(chain *certChain) error {
	err := chain.verify(s.roots)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	s.certs[fingerprintBase64(chain.leaf)] = chain

	return nil
}

func (s *certStore) get(fingerprint string) (*certChain, bool) {
	s.RLock()
	defer s.RUnlock()

	chain, ok := s.certs[fingerprint]

	return chain, ok
}

type certChain struct {
	intermediates *x509.CertPool
	leaf          *x509.Certificate
}

func parseTokenCerts(token *jwt.Token) (*certChain, error) {
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

		certs[i] = c
	}

	slices.Reverse(certs)

	lastIdx := len(certs) - 1

	intermediates := x509.NewCertPool()
	for _, cert := range certs[:lastIdx] {
		intermediates.AddCert(cert)
	}

	return &certChain{
		intermediates: intermediates,
		leaf:          certs[lastIdx],
	}, nil
}

func (c *certChain) publicKey() interface{} {
	return c.leaf.PublicKey
}

func (c *certChain) verify(roots *x509.CertPool) error {
	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: c.intermediates,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	if _, err := c.leaf.Verify(opts); err != nil {
		return err
	}

	if c.leaf.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
		return errors.New("invalid key usage")
	}

	return nil
}

func fingerprintBase64(cert *x509.Certificate) string {
	sum := sha256.Sum256(cert.Raw)
	return base64.StdEncoding.EncodeToString(sum[:])
}
