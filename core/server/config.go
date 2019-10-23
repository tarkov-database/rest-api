package server

import (
	"errors"
	"log"
	"os"
	"strconv"
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
	Port        int
	TLS         bool
	Certificate string
	PrivateKey  string
}

func newConfig() (*config, error) {
	c := &config{Port: 8080}

	if env := os.Getenv("SERVER_PORT"); len(env) > 1 {
		if i, err := strconv.Atoi(env); err == nil {
			c.Port = i
		} else {
			return c, errors.New("server port is not an integer")
		}
	}

	if env := os.Getenv("SERVER_TLS"); len(env) > 3 {
		if b, err := strconv.ParseBool(env); err == nil {
			c.TLS = b
		} else {
			return c, errors.New("invalid boolean in environment variable")
		}
	}

	if c.TLS {
		if env := os.Getenv("SERVER_CERT"); len(env) > 0 {
			c.Certificate = env
		} else {
			return c, errors.New("server certificate missing")
		}

		if env := os.Getenv("SERVER_KEY"); len(env) > 0 {
			c.PrivateKey = env
		} else {
			return c, errors.New("server private key missing")
		}
	}

	return c, nil
}
