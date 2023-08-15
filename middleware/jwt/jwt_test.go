package jwt

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/logger"
	"github.com/julienschmidt/httprouter"
)

func init() {
	logger.Init("default", false, false, io.Discard)

	// Load root certificate
	certs, err := parseCertsFromPEM("testdata/root.crt")
	if err != nil {
		log.Printf("failed to read cert file: %v", err)
		os.Exit(2)
	}

	for _, cert := range certs {
		if cert.IsCA {
			store.roots.AddCert(cert)
		}
	}
}

const badToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestSymmetricSignVerify(t *testing.T) {
	claimsIn := &Claims{}

	token, err := SignToken(claimsIn, nil)
	if err != nil {
		t.Fatalf("Token creation failed: %v", err)
	}

	claimsOut, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("Token verification failed: %v", err)
	}

	if !reflect.DeepEqual(claimsIn, claimsOut) {
		t.Fatalf("Token claim validation failed: claims %v and %v unequal", claimsIn, claimsOut)
	}

	if _, err := VerifyToken(badToken); err == nil {
		t.Error("Token verification failed: invalid token verified as valid")
	}
}

func TestAsymmetricVerify(t *testing.T) {
	tests := []struct {
		name              string
		claims            *Claims
		certFile          string
		keyFile           string
		expectedErr       error
		expectedCustomErr interface{}
	}{
		{
			name:     "valid token",
			claims:   &Claims{},
			certFile: "testdata/certs.crt",
			keyFile:  "testdata/key.pem",
		},
		{
			name:              "invalid token",
			claims:            &Claims{},
			certFile:          "testdata/certs_invalid.crt",
			keyFile:           "testdata/key_invalid.pem",
			expectedCustomErr: x509.UnknownAuthorityError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := signTestingToken(tt.claims, tt.certFile, tt.keyFile)
			if err != nil {
				t.Fatalf("Token creation failed: %v", err)
			}

			_, err = VerifyToken(token)
			if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
				t.Fatalf("Token verification failed: expected error %v, got %v", tt.expectedErr, err)
			}
			if tt.expectedCustomErr != nil && !errors.As(err, &tt.expectedCustomErr) {
				t.Fatalf("Token verification failed: expected error type %T, got %v", tt.expectedCustomErr, err)
			}
			if tt.expectedErr == nil && tt.expectedCustomErr == nil && err != nil {
				t.Errorf("Token verification failed: %v", err)
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", badToken))

	token, err := ExtractToken(&http.Request{Header: header})
	if err != nil {
		t.Fatalf("Getting token failed: %v", err)
	}

	if token != badToken {
		t.Fatalf("Getting token failed: token %s and %s unequal", token, badToken)
	}
}

func TestAuhtorizationHandler(t *testing.T) {
	tests := []struct {
		name         string
		claims       *Claims
		certFile     string
		keyFile      string
		expectedCode int
	}{
		{
			name: "valid token",
			claims: &Claims{
				Scope: []string{ScopeUserRead},
			},
			certFile:     "testdata/certs.crt",
			keyFile:      "testdata/key.pem",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid token",
			claims:       &Claims{},
			certFile:     "testdata/certs_invalid.crt",
			keyFile:      "testdata/key_invalid.pem",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid scope token",
			claims: &Claims{
				Scope: []string{ScopeItemRead},
			},
			certFile:     "testdata/certs.crt",
			keyFile:      "testdata/key.pem",
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := signTestingToken(tt.claims, tt.certFile, tt.keyFile)
			if err != nil {
				t.Fatalf("Getting token failed: %s", err)
			}

			header := http.Header{}
			header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

			handle := AuhtorizationHandler(ScopeUserRead, func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
				w.WriteHeader(http.StatusOK)
			})

			w := httptest.NewRecorder()
			handle(w, &http.Request{Header: header}, httprouter.Params{})

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Fatalf("Authorization handler failed: unexpected response code %v", resp.StatusCode)
			}
		})
	}
}

func signTestingToken(c *Claims, certPath, keyPath string) (string, error) {
	// Load key and PEM certificates
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read key file: %w", err)
	}

	key, err := jwt.ParseEdPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse key: %w", err)
	}

	certs, err := parseCertsFromPEM(certPath)
	if err != nil {
		return "", fmt.Errorf("failed to read cert file: %w", err)
	}

	// Create JWT
	now := time.Now()

	c.Audience = append(c.Audience, cfg.Audience)
	c.IssuedAt = jwt.NewNumericDate(now)
	c.ExpiresAt = jwt.NewNumericDate(now.Add(5 * time.Minute))

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, c)

	// Set x5c header with leaf and intermediate certificates
	slices.Reverse(certs)

	x5c := make([]string, len(certs))
	for i, cert := range certs {
		x5c[i] = base64.StdEncoding.EncodeToString(cert.Raw)
	}
	token.Header["x5c"] = x5c

	// Set x5t#256 header
	hash := sha256.Sum256(certs[0].Raw)
	token.Header["x5t#S256"] = base64.StdEncoding.EncodeToString(hash[:])

	// Sign JWT
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
