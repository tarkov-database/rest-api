package jwt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
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
	Algorithm      jwt.Algorithm
	Audience       jwt.Audience
	ExpirationTime time.Duration
}

func newConfig() (*config, error) {
	c := &config{}

	if env := os.Getenv("JWT_KEY"); len(env) > 0 {
		switch strings.ToLower(os.Getenv("JWT_ALG")) {
		case "hs256":
			c.Algorithm = jwt.NewHS256([]byte(env))
		case "hs512":
			c.Algorithm = jwt.NewHS512([]byte(env))
		default:
			log.Println("No valid JWT algorithm set, using HS256")
			c.Algorithm = jwt.NewHS256([]byte(env))
		}
	} else {
		return c, errors.New("jwt key is not set")
	}

	if env := os.Getenv("JWT_AUDIENCE"); len(env) >= 3 {
		c.Audience = strings.Split(env, ",")
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

	return c, nil
}
