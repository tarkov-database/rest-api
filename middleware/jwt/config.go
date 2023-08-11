package jwt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	SigningAlgorithm jwt.SigningMethod
	SigningKey       []byte
	Audience         string
	ExpirationTime   time.Duration
	Leeway           time.Duration
}

func newConfig() (*config, error) {
	c := &config{}

	if key := os.Getenv("JWT_KEY"); len(key) > 0 {
		switch strings.ToLower(os.Getenv("JWT_ALG")) {
		case "hs256":
			c.SigningAlgorithm = jwt.SigningMethodHS256
		case "hs512":
			c.SigningAlgorithm = jwt.SigningMethodHS512
		default:
			log.Println("No valid JWT algorithm set, using HS256")
			c.SigningAlgorithm = jwt.SigningMethodHS256
		}
		c.SigningKey = []byte(key)
	} else {
		return c, errors.New("jwt key is not set")
	}

	if env := os.Getenv("JWT_AUDIENCE"); len(env) >= 3 {
		c.Audience = env
	} else {
		return c, errors.New("jwt audience is not set or too short")
	}

	if env := os.Getenv("JWT_EXPIRATION"); len(env) > 0 {
		d, err := time.ParseDuration(env)
		if err != nil {
			return c, fmt.Errorf("jwt expiration value is not valid: %s", err)
		}
		c.ExpirationTime = d
	} else {
		c.ExpirationTime = 30 * time.Minute
	}

	if env := os.Getenv("JWT_LEEWAY"); len(env) > 0 {
		d, err := time.ParseDuration(env)
		if err != nil {
			return c, fmt.Errorf("jwt leeway value is not valid: %s", err)
		}
		c.Leeway = d
	} else {
		c.Leeway = 30 * time.Second
	}

	return c, nil
}
