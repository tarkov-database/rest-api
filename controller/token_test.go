package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tarkov-database/rest-api/middleware/jwt"
	"github.com/tarkov-database/rest-api/model/token"

	"github.com/julienschmidt/httprouter"
)

func TestTokenGET(t *testing.T) {
	userID := userIDs[0]

	clmIn := &jwt.Claims{}
	clmIn.Subject = userID.Hex()

	tkn, err := jwt.SignToken(clmIn, nil)
	if err != nil {
		t.Fatalf("Getting token failed: %s", err)
	}

	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	w := httptest.NewRecorder()

	TokenGET(w, &http.Request{Header: header}, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Getting token failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting token failed: content type is invalid")
	}

	output := &token.Response{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Getting token failed: %s", err)
	}

	clmOut, err := jwt.VerifyToken(output.Token)
	if err != nil {
		t.Error("Getting token failed: token invalid")
	}

	if clmOut.Subject != clmIn.Subject {
		t.Error("Getting token failed: subject invalid")
	}
}

func TestTokenPOST(t *testing.T) {
	userID := userIDs[0]

	buf := new(bytes.Buffer)

	clm := &jwt.Claims{
		Scope: []string{
			jwt.ScopeTokenWrite,
		},
	}
	clm.Subject = userID.Hex()

	tkn, err := jwt.SignToken(clm, nil)
	if err != nil {
		t.Fatalf("Creating token failed: %s", err)
	}

	input := &token.Request{
		Scope: []string{
			jwt.ScopeAllRead,
			jwt.ScopeAllWrite,
		},
		ExpiresIn: "24h",
	}
	input.Subject = userID.Hex()

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Creating token failed: %s", err)
	}

	req := httptest.NewRequest("POST", "http://example.com/v2/token", buf)
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tkn))

	w := httptest.NewRecorder()

	TokenPOST(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Creating token failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Creating token failed: content type is invalid")
	}

	output := &token.Response{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Creating token failed: %s", err)
	}

	clmOut, err := jwt.VerifyToken(output.Token)
	if err != nil {
		t.Error("Getting token failed: token invalid")
	}

	if clmOut.Subject != input.Subject {
		t.Error("Getting token failed: subject invalid")
	}

	d := clmOut.ExpiresAt.Sub(time.Now())

	if d.Round(time.Hour) != 24*time.Hour {
		t.Error("Getting token failed: lifetime mismatch")
	}
}
