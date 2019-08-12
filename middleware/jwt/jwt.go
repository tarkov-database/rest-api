package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tarkov-database/api/model"
	"github.com/tarkov-database/api/view"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/logger"
	"github.com/julienschmidt/httprouter"
)

var (
	ErrInvalidScope    = errors.New("no or invalid scopes")
	ErrInvalidSubject  = errors.New("invalid subject")
	ErrInvalidAudience = errors.New("audience mismatch")
	ErrExpiredToken    = errors.New("token is expired")
)

const (
	// Scopes
	ScopeAllRead  = "read:all"
	ScopeAllWrite = "write:all"

	ScopeItemRead  = "read:item"
	ScopeItemWrite = "write:item"

	ScopeUserRead  = "read:user"
	ScopeUserWrite = "write:user"

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

type Claims struct {
	jwt.StandardClaims `json:",inline"`
	Scope              []string `json:"scope"`
}

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

func VerifyToken(tokenStr string) (*Claims, error) {
	parser := jwt.Parser{SkipClaimsValidation: true}
	token, err := parser.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		logger.Error(err)
		return &Claims{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		logger.Error("JWT claims parsing error")
		return claims, errors.New("claims parsing error")
	}
	if err := claims.Valid(); err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			return claims, ErrExpiredToken
		}
		logger.Errorf("JWT with subject \"%v\" is invalid: %v", claims.Subject, err)
		return claims, err
	}
	if !claims.VerifyAudience(audience, true) {
		logger.Errorf("JWT with subject \"%v\" is not valid", claims.Subject)
		return claims, ErrInvalidAudience
	}

	return claims, nil
}

func CreateToken(c *Claims) (string, error) {
	c.Audience = audience
	c.IssuedAt = time.Now().Unix()
	c.ExpiresAt = time.Now().Unix() + expTime
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString(key)
}

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

func AuhtorizationHandler(scope string, h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := GetToken(r)
		if err != nil {
			statusHandler(err.Error(), http.StatusUnauthorized, w)
			return
		}

		claims, err := VerifyToken(token)
		if err != nil {
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
	res := &model.Response{}
	res.New(msg, status)
	view.RenderJSON(w, res, status)
}
