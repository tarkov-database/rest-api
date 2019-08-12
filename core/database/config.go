package database

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	URI         string
	Database    string
	TLS         bool
	Certificate string
	PrivateKey  string
	RootCA      string
}

func newConfig() *config {
	c := &config{}

	if env := os.Getenv("MONGO_URI"); len(env) > 0 {
		if !strings.HasPrefix(env, "mongodb://") && !strings.HasPrefix(env, "mongodb+srv://") {
			logger.Error("MongoDB URI not valid!")
			os.Exit(2)
		}
		c.URI = env
	} else {
		logger.Error("MongoDB URI not set!")
		os.Exit(2)
	}
	if env := os.Getenv("MONGO_DB"); len(env) > 0 {
		c.Database = env
	} else {
		logger.Error("MongoDB DB not set!")
		os.Exit(2)
	}
	if env := os.Getenv("MONGO_TLS"); len(env) > 0 {
		if b, err := strconv.ParseBool(env); err == nil {
			c.TLS = b
		}
		if env := os.Getenv("MONGO_CERT"); len(env) > 0 {
			c.Certificate = env
			if env := os.Getenv("MONGO_KEY"); len(env) > 0 {
				c.PrivateKey = env
			}
			if env := os.Getenv("MONGO_CA"); len(env) > 0 {
				c.RootCA = env
			}
		}
	}

	return c
}

func (c *config) getOptions() *options.ClientOptions {
	clientOptions := options.Client()
	clientOptions.ApplyURI(c.URI)

	if c.TLS {
		var tlsConfig *tls.Config
		var certAuth bool

		if len(c.Certificate) > 0 {
			certAuth = true
		}

		if certAuth {
			clientCert, rootCAs, err := c.getTLSCertificate()
			if err != nil {
				logger.Fatal(err)
			}

			tlsConfig = &tls.Config{
				Certificates: []tls.Certificate{clientCert},
				RootCAs:      rootCAs,
			}
		} else {
			tlsConfig = &tls.Config{}
		}

		clientOptions.SetTLSConfig(tlsConfig)
	}

	return clientOptions
}

func (c *config) getTLSCertificate() (tls.Certificate, *x509.CertPool, error) {
	rootCAs := x509.NewCertPool()

	clientCertPEM, err := ioutil.ReadFile(c.Certificate)
	if err != nil {
		return tls.Certificate{}, rootCAs, err
	}

	clientKeyPEM, err := ioutil.ReadFile(c.PrivateKey)
	if err != nil {
		return tls.Certificate{}, rootCAs, err
	}

	clientCert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
	if err != nil {
		return tls.Certificate{}, rootCAs, err
	}

	if len(c.RootCA) > 0 {
		caPEM, err := ioutil.ReadFile(c.RootCA)
		if err != nil {
			return tls.Certificate{}, rootCAs, err
		}

		ok := rootCAs.AppendCertsFromPEM(caPEM)
		if ok != true {
			return tls.Certificate{}, rootCAs, errors.New("Failed to load root CA")
		}
	}

	return clientCert, rootCAs, nil
}
