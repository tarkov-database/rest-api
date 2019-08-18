package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/tarkov-database/rest-api/core/database"
	"github.com/tarkov-database/rest-api/model/user"

	"github.com/google/logger"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const contentTypeJSON = "application/json"

var userIDs []primitive.ObjectID

func init() {
	logger.Init("default", false, false, ioutil.Discard)
}

func mongoStartup() {
	if err := database.Init(); err != nil {
		log.Fatalf("Database startup error: %s", err)
	}

	c := database.GetDB().Collection(user.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	userA := user.User{ID: createObjectID()}
	userB := user.User{ID: createObjectID()}

	if _, err := c.InsertMany(ctx, bson.A{userA, userB}); err != nil {
		log.Fatalf("Database startup error: %s", err)
	}
}

func mongoCleanup() {
	c := database.GetDB().Collection(user.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	if _, err := c.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": userIDs}}); err != nil {
		log.Fatalf("Database cleanup error: %s", err)
	}

	if err := database.Shutdown(); err != nil {
		log.Fatalf("Database shutdown error: %s", err)
	}
}

func createObjectID() primitive.ObjectID {
	userID := primitive.NewObjectID()
	userIDs = append(userIDs, userID)

	return userID
}

func removeObjectID(id primitive.ObjectID) {
	new := make([]primitive.ObjectID, 0, len(userIDs)-1)
	for _, k := range userIDs {
		if k != id {
			new = append(new, k)
		}
	}

	userIDs = new
}

func TestMain(m *testing.M) {
	mongoStartup()
	code := m.Run()
	mongoCleanup()
	os.Exit(code)
}

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

type result struct {
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

	res := &result{}

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
	userID := createObjectID()

	buf := new(bytes.Buffer)

	newUser := &user.User{ID: userID, Email: "test@testing.dev"}

	if err := json.NewEncoder(buf).Encode(newUser); err != nil {
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

	usr := &user.User{}

	if err := json.NewDecoder(resp.Body).Decode(usr); err != nil {
		t.Fatalf("Creating user failed: %s", err)
	}

	if usr.ID != userID {
		t.Errorf("Creating user failed: user ID %s and %s unequal", usr.ID, userID)
	}
}

func TestUserPUT(t *testing.T) {
	userID := userIDs[0]

	email := "test2@testing.dev"

	buf := new(bytes.Buffer)

	newData := &user.User{Email: email}

	if err := json.NewEncoder(buf).Encode(newData); err != nil {
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

	usr := &user.User{}

	if err := json.NewDecoder(resp.Body).Decode(usr); err != nil {
		t.Fatalf("Replacing user failed: %s", err)
	}

	if usr.Email != email {
		t.Errorf("Replacing user failed: user eMail %s and %s unequal", usr.Email, email)
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

	removeObjectID(userID)
}
