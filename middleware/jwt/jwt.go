package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/view"

	"github.com/dgrijalva/jwt-go"
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

const (
	// ScopeAllRead represents the global read permission scope
	ScopeAllRead = "read:all"

	// ScopeAllWrite represents the global write permission scope
	ScopeAllWrite = "write:all"

	// ScopeItemRead represents the item read permission scope
	ScopeItemRead = "read:item"

	// ScopeItemWrite represents the item write permission scope
	ScopeItemWrite = "write:item"

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
	case ScopeUserRead:
		valid = true
	case ScopeUserWrite:
		valid = true
	case ScopeTokenWrite:
		valid = true
	}

	return valid
}

var (
	key      []byte
	audience string
	expTime  int64 = 30 * 60
)

func init() {
	if env := os.Getenv("JWT_KEY"); len(env) >= 8 {
		key = []byte(env)
	} else {
		logger.Fatal(errors.New("JWT key is not set or too short"))
	}

	if env := os.Getenv("JWT_AUDIENCE"); len(env) >= 3 {
		audience = env
	} else {
		logger.Fatal(errors.New("JWT audience is not set or too short"))
	}

	if env := os.Getenv("JWT_EXPIRATION"); len(env) > 0 {
		if s, err := strconv.ParseInt(env, 10, 64); err == nil {
			expTime = s * 60
		}
	}
}

// Claims represents the claims of a token
type Claims struct {
	jwt.StandardClaims `json:",inline"`
	Scope              []string `json:"scope"`
}

// ValidateCustom validates the custom claims of a token
func (c *Claims) ValidateCustom() error {
	if len(c.Subject) < 24 {
		return ErrInvalidSubject
	}
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
	c.Audience = audience
	c.IssuedAt = time.Now().Unix()
	c.ExpiresAt = time.Now().Unix() + expTime
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString(key)
}

// GetToken gets the token of an HTTP request
func GetToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if len(header) == 0 {
		return "", errors.New("authorization header not set")
	}

	headerString := strings.TrimSpace(header)
	if !strings.HasPrefix(header, "Bearer") {
		return "", errors.New("invalid authorization header")
	}

	return strings.TrimSpace(strings.TrimPrefix(headerString, "Bearer")), nil
}

// VerifyToken verifies a token
func VerifyToken(tokenStr string) (*Claims, error) {
	parser := jwt.Parser{SkipClaimsValidation: true}
	token, err := parser.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return &Claims{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return claims, errors.New("claims parsing error")
	}
	if err := claims.Valid(); err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			return claims, ErrExpiredToken
		}
		return claims, err
	}
	if !claims.VerifyAudience(audience, true) {
		return claims, ErrInvalidAudience
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
			if err != ErrExpiredToken {
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
