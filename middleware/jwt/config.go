package jwt

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

var cfg *config

func init() {
	cfg = newConfig()
}

type config struct {
	Algorithm      jwt.Algorithm
	Audience       jwt.Audience
	ExpirationTime time.Duration
}

func newConfig() *config {
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
		log.Fatal("JWT key is not set or too short")
	}

	if env := os.Getenv("JWT_AUDIENCE"); len(env) >= 3 {
		c.Audience = strings.Split(env, ",")
	} else {
		log.Fatal("JWT audience is not set or too short")
	}

	if env := os.Getenv("JWT_EXPIRATION"); len(env) > 0 {
		d, err := time.ParseDuration(env)
		if err != nil {
			log.Fatalf("JWT expiration value is not valid: %s", err)
		}
		c.ExpirationTime = d
	}

	return c
}
