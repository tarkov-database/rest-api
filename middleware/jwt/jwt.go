package jwt

import (
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/view"

	"github.com/golang-jwt/jwt/v5"
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

	// ErrNotBefore indicates that the token is not yet valid
	ErrNotBefore = errors.New("token is not yet valid")

	// ErrMalformed indicates that the token is malformed
	ErrMalformed = errors.New("token is malformed")

	// ErrInvalidToken indicates that the token is invalid
	ErrInvalidToken = errors.New("token is invalid")
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

	// ScopeHideoutRead represents the hideout read permission scope
	ScopeHideoutRead = "read:hideout"

	// ScopeHideoutWrite represents the hideout write permission scope
	ScopeHideoutWrite = "write:hideout"

	// ScopeLocationRead represents the location read permission scope
	ScopeLocationRead = "read:location"

	// ScopeLocationWrite represents the location write permission scope
	ScopeLocationWrite = "write:location"

	// ScopeStatisticRead represents the statistic read permission scope
	ScopeStatisticRead = "read:statistic"

	// ScopeStatisticWrite represents the statistic write permission scope
	ScopeStatisticWrite = "write:statistic"

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
	case ScopeHideoutRead:
		valid = true
	case ScopeHideoutWrite:
		valid = true
	case ScopeLocationRead:
		valid = true
	case ScopeLocationWrite:
		valid = true
	case ScopeStatisticRead:
		valid = true
	case ScopeStatisticWrite:
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
	jwt.RegisteredClaims
	Scope []string `json:"scope"`
}

// ValidateCustom validates the custom claims of a token
func (c *Claims) Validate() error {
	for _, s := range c.Scope {
		if !isScopeValid(s) {
			return ErrInvalidScope
		}
	}

	return nil
}

// SignToken signs a token
func SignToken(c *Claims, d *time.Duration) (string, error) {
	now := time.Now()

	c.Audience = append(c.Audience, cfg.Audience)
	c.IssuedAt = jwt.NewNumericDate(now)

	if d != nil {
		c.ExpiresAt = jwt.NewNumericDate(now.Add(*d))
	} else {
		c.ExpiresAt = jwt.NewNumericDate(now.Add(cfg.ExpirationTime))
	}

	token := jwt.NewWithClaims(cfg.SigningAlgorithm, c)
	s, err := token.SignedString(cfg.SigningKey)
	if err != nil {
		return "", err
	}

	return s, nil
}

// ExtractToken extracts a token from a request header
func ExtractToken(r *http.Request) (string, error) {
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

	leeway := jwt.WithLeeway(cfg.Leeway)
	audience := jwt.WithAudience(cfg.Audience)

	_, err := jwt.ParseWithClaims(tokenStr, claims, keyFunc, leeway, audience)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, errors.Join(ErrExpiredToken, err)
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, errors.Join(ErrNotBefore, err)
		case errors.Is(err, jwt.ErrTokenInvalidAudience):
			return nil, errors.Join(ErrInvalidAudience, err)
		case errors.Is(err, jwt.ErrTokenInvalidSubject):
			return nil, errors.Join(ErrInvalidSubject, err)
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, errors.Join(ErrMalformed, err)
		case errors.Is(err, jwt.ErrInvalidKey), errors.Is(err, jwt.ErrInvalidKeyType):
			return nil, err
		default:
			return nil, errors.Join(ErrInvalidToken, err)
		}
	}

	return claims, nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	switch alg := token.Method.Alg(); alg {
	// RSA algorithms
	case jwt.SigningMethodRS256.Alg(), jwt.SigningMethodRS384.Alg(), jwt.SigningMethodRS512.Alg():
	// RSAPSS algorithms
	case jwt.SigningMethodPS256.Alg(), jwt.SigningMethodPS384.Alg(), jwt.SigningMethodPS512.Alg():
	// ECDSA algorithms
	case jwt.SigningMethodES256.Alg(), jwt.SigningMethodES384.Alg(), jwt.SigningMethodES512.Alg():
	// EdDSA algorithms
	case jwt.SigningMethodEdDSA.Alg():
	// HMAC algorithms
	case jwt.SigningMethodHS256.Alg(), jwt.SigningMethodHS384.Alg(), jwt.SigningMethodHS512.Alg():
		return cfg.SigningKey, nil
	default:
		return nil, fmt.Errorf("unsupported signing algorithm: %s", alg)
	}

	fingerprint, ok := token.Header["x5t#S256"].(string)
	if !ok {
		return nil, errors.New("invalid fingerprint")
	}

	cert, ok := store.get(fingerprint)
	if ok {
		if time.Now().After(cert.NotAfter) {
			store.remove(fingerprint)
			return nil, errors.New("certificate expired")
		}
	} else {
		certs, err := parseCertsFromToken(token)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate chain: %w", err)
		}

		leaf := certs[0]

		var intermediates []*x509.Certificate
		if len(certs) > 1 {
			intermediates = certs[1:]
		}

		if err := verifyCert(leaf, intermediates, store.roots); err != nil {
			return nil, fmt.Errorf("failed to verify certificate: %w", err)
		}

		store.add(leaf)
		cert = leaf
	}

	return cert.PublicKey, nil
}

// AuhtorizationHandler returns a JWT authorization handler
func AuhtorizationHandler(scope string, h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := ExtractToken(r)
		if err != nil {
			AddAuthenticateHeader(w, err, scope)
			statusHandler(err.Error(), http.StatusUnauthorized, w)
			return
		}

		var allScope string
		if scope != "" {
			allScope = fmt.Sprintf("%s:all", strings.SplitN(scope, ":", 2)[0])
		}

		claims, err := VerifyToken(token)
		if err != nil {
			AddAuthenticateHeader(w, err, scope, allScope)
			statusHandler(err.Error(), http.StatusUnauthorized, w)
			return
		}

		var ok bool
		if scope == "" {
			ok = true
		} else {
			for _, s := range claims.Scope {
				if s == scope || s == allScope {
					ok = true
					break
				}
			}
		}

		if !ok {
			AddAuthenticateHeader(w, ErrInvalidScope, scope, allScope)
			statusHandler("Insufficient permissions", http.StatusForbidden, w)
			return
		}

		h(w, r, ps)
	}
}

const (
	authenticateInvalid      = "invalid_token"
	authenticateInsufficient = "insufficient_scope"
)

// AddAuthenticateHeader adds the "WWW-Authenticate" header to the response
func AddAuthenticateHeader(w http.ResponseWriter, err error, scopes ...string) {
	value := fmt.Sprintf("Bearer scope=\"%s\"", strings.Join(scopes, " "))

	switch err {
	case ErrExpiredToken, ErrNotBefore, ErrInvalidAudience, ErrInvalidSubject, ErrMalformed, ErrInvalidToken:
		value += fmt.Sprintf(", error=\"%s\"", authenticateInvalid)
	case ErrInvalidScope:
		value += fmt.Sprintf(", error=\"%s\"", authenticateInsufficient)
	}

	w.Header().Add("WWW-Authenticate", value)
}

func statusHandler(msg string, status int, w http.ResponseWriter) {
	res := model.NewResponse(msg, status)
	view.RenderJSON(res, status, w)
}
