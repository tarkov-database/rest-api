package jwt

import (
	"crypto/sha256"
	"encoding/base64"
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

func TestHmacSignVerify(t *testing.T) {
	claimsIn := &Claims{}

	token, err := SignToken(claimsIn, nil)
	if err != nil {
		t.Errorf("Token creation failed: %v", err)
	}

	claimsOut, err := VerifyToken(token)
	if err != nil {
		t.Errorf("Token verification failed: %v", err)
	}

	if !reflect.DeepEqual(claimsIn, claimsOut) {
		t.Errorf("Token claim validation failed: claims %v and %v unequal", claimsIn, claimsOut)
	}

	if _, err := VerifyToken(badToken); err == nil {
		t.Error("Token verification failed: invalid token verified as valid")
	}
}

func TestEd25519Verify(t *testing.T) {
	claimsIn := &Claims{}

	token, err := signTestingToken(claimsIn)
	if err != nil {
		t.Errorf("Token creation failed: %v", err)
	}

	claimsOut, err := VerifyToken(token)
	if err != nil {
		t.Errorf("Token verification failed: %v", err)
	}

	if !reflect.DeepEqual(claimsIn, claimsOut) {
		t.Errorf("Token claim validation failed: claims %v and %v unequal", claimsIn, claimsOut)
	}
}

func TestGetToken(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", badToken))

	token, err := ExtractToken(&http.Request{Header: header})
	if err != nil {
		t.Errorf("Getting token failed: %v", err)
	}

	if token != badToken {
		t.Errorf("Getting token failed: token %s and %s unequal", token, badToken)
	}
}

func TestAuhtorizationHandler(t *testing.T) {
	validClaims := &Claims{
		Scope: []string{
			ScopeUserRead,
		},
	}

	validToken, err := signTestingToken(validClaims)
	if err != nil {
		t.Fatalf("Getting token failed: %s", err)
	}

	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", validToken))

	handle := AuhtorizationHandler(ScopeUserRead, func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
	})

	w := httptest.NewRecorder()

	handle(w, &http.Request{Header: header}, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Authorization handler failed: unexpcted response code %v", resp.StatusCode)
	}

	header.Set("Authorization", fmt.Sprintf("Bearer %s", badToken))

	w = httptest.NewRecorder()

	handle(w, &http.Request{Header: header}, httprouter.Params{})

	resp = w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Authorization handler failed: unexpcted response code %v", resp.StatusCode)
	}

	invalidClaims := &Claims{
		Scope: []string{
			ScopeItemRead,
		},
	}

	invalidToken, err := signTestingToken(invalidClaims)
	if err != nil {
		t.Fatalf("Authorization handler failed: %s", err)
	}

	header.Set("Authorization", fmt.Sprintf("Bearer %s", invalidToken))

	w = httptest.NewRecorder()

	handle(w, &http.Request{Header: header}, httprouter.Params{})

	resp = w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Authorization handler failed: unexpcted response code %v", resp.StatusCode)
	}
}

func signTestingToken(c *Claims) (string, error) {
	// Load key and PEM certificates
	keyBytes, err := os.ReadFile("testdata/key.pem")
	if err != nil {
		return "", fmt.Errorf("failed to read key file: %w", err)
	}

	key, err := jwt.ParseEdPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse key: %w", err)
	}

	certs, err := parseCertsFromPEM("testdata/certs.crt")
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
