package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tarkov-database/rest-api/model/location"

	"github.com/julienschmidt/httprouter"
)

type locationResult struct {
	Count int64               `json:"total"`
	Items []location.Location `json:"items"`
}

func TestLocationGET(t *testing.T) {
	locationID := locationIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: locationID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	LocationGET(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting location failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting location failed: content type is invalid")
	}

	output := &location.Location{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Getting location failed: %s", err)
	}

	if output.ID != locationID {
		t.Error("Getting location failed: location ID invalid")
	}
}

func TestLocationsGET(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/v2/location", nil)

	w := httptest.NewRecorder()

	LocationsGET(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting locations failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting locations failed: content type is invalid")
	}

	res := &locationResult{}

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		t.Fatalf("Getting locations failed: %s", err)
	}

	if res.Count < 1 {
		t.Error("Getting locations failed: result count invalid")
	}
	if len(res.Items) == 0 {
		t.Fatal("Getting locations failed: result empty")
	}
}

func TestLocationPOST(t *testing.T) {
	locationID := createLocationID()

	buf := new(bytes.Buffer)

	input := &location.Location{
		ID:             locationID,
		Name:           "location",
		Description:    "a test location",
		MinimumPlayers: 2,
		MaximumPlayers: 6,
		Available:      true,
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Creating location failed: %s", err)
	}

	req := httptest.NewRequest("POST", "http://example.com/v2/location", buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	w := httptest.NewRecorder()

	LocationPOST(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Creating location failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Creating location failed: content type is invalid")
	}

	output := &location.Location{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Creating location failed: %s", err)
	}

	if output.ID != input.ID {
		t.Errorf("Creating location failed: location ID %s and %s unequal", output.ID, input.ID)
	}
}

func TestLocationPUT(t *testing.T) {
	locationID := locationIDs[0]

	buf := new(bytes.Buffer)

	input := &location.Location{
		ID:             locationID,
		Name:           "Location",
		Description:    "A test location",
		MinimumPlayers: 2,
		MaximumPlayers: 6,
		Available:      false,
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Replacing location failed: %s", err)
	}

	req := httptest.NewRequest("PUT", "http://example.com/v2/location", buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: locationID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	LocationPUT(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Replacing location failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Replacing location failed: content type is invalid")
	}

	output := &location.Location{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Replacing location failed: %s", err)
	}

	if output.Name != input.Name {
		t.Errorf("Replacing location failed: location e-mail %s and %s unequal", output.Name, input.Name)
	}
}

func TestLocationDELETE(t *testing.T) {
	locationID := locationIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: locationID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	LocationDELETE(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Deleting location failed: unexpcted response code %v", resp.StatusCode)
	}

	removeLocationID(locationID)
}
