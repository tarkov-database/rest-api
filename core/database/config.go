package database

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/mongo/options"
)

var cfg *config

func init() {
	var err error

	cfg, err = newConfig()
	if err != nil {
		log.Printf("Configuration error: %s\n", err)
		os.Exit(2)
	}
}

type config struct {
	URI         string
	Database    string
	TLS         bool
	Certificate string
	PrivateKey  string
	RootCA      string
}

func newConfig() (*config, error) {
	c := &config{}

	if env := os.Getenv("MONGO_URI"); len(env) > 0 {
		if !strings.HasPrefix(env, "mongodb://") && !strings.HasPrefix(env, "mongodb+srv://") {
			return c, errors.New("mongo uri invalid")
		}
		c.URI = env
	} else {
		return c, errors.New("mongo uri not set")
	}

	if env := os.Getenv("MONGO_DB"); len(env) > 0 {
		c.Database = env
	} else {
		return c, errors.New("mongo database not set")
	}

	if env := os.Getenv("MONGO_TLS"); len(env) > 0 {
		if b, err := strconv.ParseBool(env); err == nil {
			c.TLS = b
		} else {
			return c, errors.New("invalid boolean in environment variable")
		}

		if env := os.Getenv("MONGO_CERT"); len(env) > 0 {
			c.Certificate = env

			if env := os.Getenv("MONGO_KEY"); len(env) > 0 {
				c.PrivateKey = env
			} else {
				return c, errors.New("mongo database not set")
			}

			if env := os.Getenv("MONGO_CA"); len(env) > 0 {
				c.RootCA = env
			}
		}
	}

	return c, nil
}

func (c *config) getClientOptions() (*options.ClientOptions, error) {
	opts := options.Client()
	opts.ApplyURI(c.URI)

	if c.TLS {
		var tlsConfig *tls.Config
		var certAuth bool

		if len(c.Certificate) > 0 {
			certAuth = true
		}

		if certAuth {
			clientCert, rootCAs, err := c.getTLSCertificate()
			if err != nil {
				return opts, fmt.Errorf("certificate loading error: %s", err)
			}

			tlsConfig = &tls.Config{
				Certificates: []tls.Certificate{clientCert},
				RootCAs:      rootCAs,
			}
		} else {
			tlsConfig = &tls.Config{}
		}

		opts.SetTLSConfig(tlsConfig)
	}

	return opts, nil
}

func (c *config) getTLSCertificate() (tls.Certificate, *x509.CertPool, error) {
	var clientCert tls.Certificate
	var rootCAs *x509.CertPool

	var err error

	clientCertPEM, err := ioutil.ReadFile(c.Certificate)
	if err != nil {
		return clientCert, rootCAs, err
	}

	clientKeyPEM, err := ioutil.ReadFile(c.PrivateKey)
	if err != nil {
		return clientCert, rootCAs, err
	}

	clientCert, err = tls.X509KeyPair(clientCertPEM, clientKeyPEM)
	if err != nil {
		return clientCert, rootCAs, err
	}

	if len(c.RootCA) > 0 {
		rootCAs = x509.NewCertPool()

		caPEM, err := ioutil.ReadFile(c.RootCA)
		if err != nil {
			return clientCert, rootCAs, err
		}

		ok := rootCAs.AppendCertsFromPEM(caPEM)
		if ok != true {
			return clientCert, rootCAs, errors.New("failed to load root CA")
		}
	}

	return clientCert, rootCAs, nil
}
