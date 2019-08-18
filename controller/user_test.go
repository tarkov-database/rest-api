package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tarkov-database/rest-api/model/user"

	"github.com/julienschmidt/httprouter"
)

func TestUserGET(t *testing.T) {
	userID := userIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: userID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	UserGET(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting user failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting user failed: content type is invalid")
	}

	usr := &user.User{}

	if err := json.NewDecoder(resp.Body).Decode(usr); err != nil {
		t.Fatalf("Getting user failed: %s", err)
	}

	if usr.ID != userID {
		t.Error("Getting user failed: user ID invalid")
	}
}

type userResult struct {
	Count int64       `json:"total"`
	Items []user.User `json:"items"`
}

func TestUsersGET(t *testing.T) {
	userID := userIDs[0]

	req := httptest.NewRequest("GET", "http://example.com/v2/user", nil)

	w := httptest.NewRecorder()

	UsersGET(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting users failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting users failed: content type is invalid")
	}

	res := &userResult{}

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		t.Fatalf("Getting users failed: %s", err)
	}

	if res.Count < 1 {
		t.Error("Getting users failed: result count invalid")
	}
	if len(res.Items) == 0 {
		t.Fatal("Getting users failed: result empty")
	}
	if id := res.Items[0].ID; id != userID {
		t.Errorf("Getting users failed: user ID %s and %s unequal", id.Hex(), userID.Hex())
	}
}

func TestUserPOST(t *testing.T) {
	userID := createUserID()

	buf := new(bytes.Buffer)

	input := &user.User{ID: userID, Email: "test@testing.dev"}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Creating user failed: %s", err)
	}

	req := httptest.NewRequest("POST", "http://example.com/v2/user", buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	w := httptest.NewRecorder()

	UserPOST(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Creating user failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Creating user failed: content type is invalid")
	}

	output := &user.User{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Creating user failed: %s", err)
	}

	if output.ID != input.ID {
		t.Errorf("Creating user failed: user ID %s and %s unequal", output.ID, input.ID)
	}
}

func TestUserPUT(t *testing.T) {
	userID := userIDs[0]

	buf := new(bytes.Buffer)

	input := &user.User{Email: "test2@testing.dev"}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Replacing user failed: %s", err)
	}

	req := httptest.NewRequest("PUT", "http://example.com/v2/user", buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: userID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	UserPUT(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Replacing user failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Replacing user failed: content type is invalid")
	}

	output := &user.User{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Replacing user failed: %s", err)
	}

	if output.Email != input.Email {
		t.Errorf("Replacing user failed: user e-mail %s and %s unequal", output.Email, input.Email)
	}
}

func TestUserDELETE(t *testing.T) {
	userID := userIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: userID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	UserDELETE(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Deleting user failed: unexpcted response code %v", resp.StatusCode)
	}

	removeUserID(userID)
}
