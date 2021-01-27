package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/view"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/google/logger"
	"github.com/julienschmidt/httprouter"
)

var (
	// ErrInvalidScope indicates that a scope value is not valid
	ErrInvalidScope = errors.New("no or invalid scopes")

	// ErrInvalidSubject indicates that the subject value is not valid
	ErrInvalidSubject = errors.New("invalid subject")

	// ErrInvalidAudience indicates that the audience value does not match
	ErrInvalidAudience = errors.New("audience mismatch")

	// ErrExpiredToken indicates that the token is expired
	ErrExpiredToken = errors.New("token is expired")
)

var (
	// ErrNoAuthHeader indicates that the authorization is not set
	ErrNoAuthHeader = errors.New("authorization header not set")

	// ErrInvalidAuthHeader indicates that the authorization is invalid
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
)

const (
	// ScopeAllRead represents the global read permission scope
	ScopeAllRead = "read:all"

	// ScopeAllWrite represents the global write permission scope
	ScopeAllWrite = "write:all"

	// ScopeItemRead represents the item read permission scope
	ScopeItemRead = "read:item"

	// ScopeItemWrite represents the item write permission scope
	ScopeItemWrite = "write:item"

	// ScopeLocationRead represents the location read permission scope
	ScopeLocationRead = "read:location"

	// ScopeLocationWrite represents the location write permission scope
	ScopeLocationWrite = "write:location"

	// ScopeUserRead represents the user read permission scope
	ScopeUserRead = "read:user"

	// ScopeUserWrite represents the user write permission scope
	ScopeUserWrite = "write:user"

	// ScopeTokenWrite represents the token write permission scope
	ScopeTokenWrite = "write:token"
)

func isScopeValid(s string) bool {
	var valid bool

	switch s {
	case ScopeAllRead:
		valid = true
	case ScopeAllWrite:
		valid = true
	case ScopeItemRead:
		valid = true
	case ScopeItemWrite:
		valid = true
	case ScopeLocationRead:
		valid = true
	case ScopeLocationWrite:
		valid = true
	case ScopeUserRead:
		valid = true
	case ScopeUserWrite:
		valid = true
	case ScopeTokenWrite:
		valid = true
	}

	return valid
}

// Claims represents the claims of a token
type Claims struct {
	jwt.Payload
	Scope []string `json:"scope"`
}

// ValidateCustom validates the custom claims of a token
func (c *Claims) ValidateCustom() error {
	if len(c.Scope) == 0 {
		return ErrInvalidScope
	}

	for _, s := range c.Scope {
		if !isScopeValid(s) {
			return ErrInvalidScope
		}
	}

	return nil
}

// CreateToken creates a new token
func CreateToken(c *Claims) (string, error) {
	now := time.Now()

	c.Audience = cfg.Audience
	c.IssuedAt = jwt.NumericDate(now)
	c.ExpirationTime = jwt.NumericDate(now.Add(cfg.ExpirationTime))

	t, err := jwt.Sign(c, cfg.Algorithm)
	if err != nil {
		return "", err
	}

	return string(t), nil
}

// GetToken gets the token of an HTTP request
func GetToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if len(header) == 0 {
		return "", ErrNoAuthHeader
	}

	headerStr := strings.TrimSpace(header)
	if !strings.HasPrefix(header, "Bearer ") {
		return "", ErrInvalidAuthHeader
	}

	return strings.TrimPrefix(headerStr, "Bearer "), nil
}

// VerifyToken verifies a token
func VerifyToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	now := time.Now()

	expVal := jwt.ExpirationTimeValidator(now)
	audVal := jwt.AudienceValidator(cfg.Audience)

	valPayload := jwt.ValidatePayload(&claims.Payload, expVal, audVal)

	if _, err := jwt.Verify([]byte(tokenStr), cfg.Algorithm, claims, valPayload); err != nil {
		switch err {
		case jwt.ErrExpValidation:
			return claims, ErrExpiredToken
		case jwt.ErrAudValidation:
			return claims, ErrInvalidAudience
		default:
			return claims, err
		}
	}

	return claims, nil
}

// AuhtorizationHandler returns a JWT authorization handler
func AuhtorizationHandler(scope string, h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := GetToken(r)
		if err != nil {
			statusHandler(err.Error(), http.StatusUnauthorized, w)
			return
		}

		claims, err := VerifyToken(token)
		if err != nil {
			if !errors.Is(err, ErrExpiredToken) {
				logger.Error(err)
			}
			statusHandler(err.Error(), http.StatusUnauthorized, w)
			return
		}

		var ok bool
		all := fmt.Sprintf("%s:all", strings.Split(scope, ":")[0])
		for _, s := range claims.Scope {
			if s == scope || s == all {
				ok = true
				break
			}
		}

		if !ok {
			statusHandler("Insufficient permissions", http.StatusForbidden, w)
			return
		}

		h(w, r, ps)
	}
}

func statusHandler(msg string, status int, w http.ResponseWriter) {
	res := model.NewResponse(msg, status)
	view.RenderJSON(res, status, w)
}
