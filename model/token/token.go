package token

import (
	"time"

	"github.com/tarkov-database/rest-api/middleware/jwt"
)

// Request represents the body of a token creation request
type Request struct {
	Subject   string   `json:"sub"`
	Scope     []string `json:"scope"`
	ExpiresIn string   `json:"expiresIn"`
}

// Duration parses ExpiresIn and returns it as time.Duration
func (r *Request) Duration() (*time.Duration, error) {
	if r.ExpiresIn == "" {
		return nil, nil
	}

	d, err := time.ParseDuration(r.ExpiresIn)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ToClaims returns jwt.Claims based on the request data
func (r *Request) ToClaims() *jwt.Claims {
	c := &jwt.Claims{}

	c.Subject = r.Subject
	c.Scope = r.Scope

	return c
}

// Response represents the body of a token creation response
type Response struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}
