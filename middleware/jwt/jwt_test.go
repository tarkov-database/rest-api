package jwt

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/logger"
	"github.com/julienschmidt/httprouter"
)

func init() {
	logger.Init("default", false, false, io.Discard)
}

const badToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestTokenCreateVerify(t *testing.T) {
	claimsIn := &Claims{}

	token, err := CreateToken(claimsIn)
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

func TestGetToken(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", badToken))

	token, err := GetToken(&http.Request{Header: header})
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

	validToken, err := CreateToken(validClaims)
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

	invalidToken, err := CreateToken(invalidClaims)
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
