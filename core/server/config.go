package server

import (
	"os"
	"strconv"
	// "github.com/google/logger"
)

type config struct {
	Port        int
	TLS         bool
	Certificate string
	PrivateKey  string
}

func (c *config) Get() {
	c.Port = 8080

	if env := os.Getenv("SERVER_PORT"); len(env) > 1 {
		if i, err := strconv.Atoi(env); err == nil {
			c.Port = int(i)
		}
	}
	if env := os.Getenv("SERVER_TLS"); len(env) > 3 {
		if b, err := strconv.ParseBool(env); err == nil {
			c.TLS = b
		}

		if env := os.Getenv("SERVER_CERT"); len(env) > 4 {
			c.Certificate = env
		}
		if env := os.Getenv("SERVER_KEY"); len(env) > 4 {
			c.PrivateKey = env
		}
	}
}
