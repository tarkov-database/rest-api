package jwt

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

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

	_, err = VerifyToken(testToken)
	if err == nil {
		t.Error("Token verification failed: invalid token verified as valid")
	}
}

func TestGetToken(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", testToken))

	token, err := GetToken(&http.Request{Header: header})
	if err != nil {
		t.Errorf("Getting token failed: %v", err)
	}

	if token != testToken {
		t.Errorf("Getting token failed: token %s and %s unequal", token, testToken)
	}
}
