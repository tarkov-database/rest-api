package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tarkov-database/rest-api/model/location"
	"github.com/tarkov-database/rest-api/model/location/feature"
	"github.com/tarkov-database/rest-api/model/location/featuregroup"

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
		t.Errorf("Replacing location failed: location name %s and %s unequal", output.Name, input.Name)
	}
}

func TestLocationDELETE(t *testing.T) {
	locationID := locationIDs[len(locationIDs)-1]

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

type featureResult struct {
	Count int64             `json:"total"`
	Items []feature.Feature `json:"items"`
}

func TestFeatureGET(t *testing.T) {
	locationID := locationIDs[0]
	featureID := featureIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: locationID.Hex(),
		},
		httprouter.Param{
			Key:   "fid",
			Value: featureID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	FeatureGET(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting feature failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting feature failed: content type is invalid")
	}

	output := &feature.Feature{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Getting feature failed: %s", err)
	}

	if output.ID != featureID {
		t.Error("Getting feature failed: feature ID invalid")
	}
}

func TestFeaturesGET(t *testing.T) {
	locationID := locationIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: locationID.Hex(),
		},
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("http://example.com/v2/location/%s/feature", locationID), nil)

	w := httptest.NewRecorder()

	FeaturesGET(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting locations failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting locations failed: content type is invalid")
	}

	res := &featureResult{}

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

func TestFeaturePOST(t *testing.T) {
	locationID := locationIDs[0]
	featureID := createFeatureID()

	buf := new(bytes.Buffer)

	input := &feature.Feature{
		ID:    featureID,
		Name:  "feature",
		Group: featureGroupIDs[0],
		Geometry: feature.Geometry{
			Type:        feature.Point,
			Coordinates: createFeatureCoords(),
		},
		Location: locationID,
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Creating feature failed: %s", err)
	}

	req := httptest.NewRequest("POST", fmt.Sprintf("http://example.com/v2/location/%s/feature", locationID), buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: locationID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	FeaturePOST(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Creating feature failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Creating feature failed: content type is invalid")
	}

	output := &feature.Feature{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Creating feature failed: %s", err)
	}

	if output.ID != input.ID {
		t.Errorf("Creating feature failed: feature ID %s and %s unequal", output.ID, input.ID)
	}
}

func TestFeaturePUT(t *testing.T) {
	locationID := locationIDs[0]
	featureID := featureIDs[0]

	buf := new(bytes.Buffer)

	input := &feature.Feature{
		ID:    featureID,
		Name:  "Feature",
		Group: featureGroupIDs[0],
		Geometry: feature.Geometry{
			Type:        feature.Point,
			Coordinates: createFeatureCoords(),
		},
		Location: locationID,
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Replacing feature failed: %s", err)
	}

	req := httptest.NewRequest("PUT", fmt.Sprintf("http://example.com/v2/location/%s/feature/%s", locationID, featureID), buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: locationID.Hex(),
		},
		httprouter.Param{
			Key:   "fid",
			Value: featureID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	FeaturePUT(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Replacing feature failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Replacing feature failed: content type is invalid")
	}

	output := &feature.Feature{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Replacing feature failed: %s", err)
	}

	if output.Name != input.Name {
		t.Errorf("Replacing feature failed: feature name %s and %s unequal", output.Name, input.Name)
	}
}

func TestFeatureDELETE(t *testing.T) {
	featureID := featureIDs[len(featureGroupIDs)-1]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: featureID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	FeatureDELETE(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Deleting location failed: unexpcted response code %v", resp.StatusCode)
	}

	removeFeatureID(featureID)
}

type featureGroupResult struct {
	Count int64                `json:"total"`
	Items []featuregroup.Group `json:"items"`
}

func TestFeatureGroupGET(t *testing.T) {
	locationID := locationIDs[0]
	featureGroupID := featureGroupIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: locationID.Hex(),
		},
		httprouter.Param{
			Key:   "gid",
			Value: featureGroupID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	FeatureGroupGET(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting feature group failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting feature group failed: content type is invalid")
	}

	output := &featuregroup.Group{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Getting feature group failed: %s", err)
	}

	if output.ID != featureGroupID {
		t.Error("Getting feature group failed: feature group ID invalid")
	}
}

func TestFeatureGroupsGET(t *testing.T) {
	locationID := locationIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: locationID.Hex(),
		},
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("http://example.com/v2/location/%s/featuregroup", locationID), nil)

	w := httptest.NewRecorder()

	FeatureGroupsGET(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting locations failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting locations failed: content type is invalid")
	}

	res := &featureGroupResult{}

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

func TestFeatureGroupPOST(t *testing.T) {
	locationID := locationIDs[0]
	featureGroupID := createFeatureGroupID()

	buf := new(bytes.Buffer)

	input := &featuregroup.Group{
		ID:          featureGroupID,
		Name:        "group",
		Description: "description",
		Tags:        []string{"test"},
		Location:    locationID,
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Creating feature group failed: %s", err)
	}

	req := httptest.NewRequest("POST", fmt.Sprintf("http://example.com/v2/location/%s/featuregroup", locationID), buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: locationID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	FeatureGroupPOST(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Creating feature group failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Creating feature group failed: content type is invalid")
	}

	output := &featuregroup.Group{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Creating feature group failed: %s", err)
	}

	if output.ID != input.ID {
		t.Errorf("Creating feature group failed: feature group ID %s and %s unequal", output.ID, input.ID)
	}
}

func TestFeatureGroupPUT(t *testing.T) {
	locationID := locationIDs[0]
	featureGroupID := featureGroupIDs[0]

	buf := new(bytes.Buffer)

	input := &featuregroup.Group{
		ID:          featureGroupID,
		Name:        "Group",
		Description: "description",
		Tags:        []string{"test"},
		Location:    locationID,
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Replacing feature group failed: %s", err)
	}

	req := httptest.NewRequest("PUT", fmt.Sprintf("http://example.com/v2/location/%s/feature/%s", locationID, featureGroupID), buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: locationID.Hex(),
		},
		httprouter.Param{
			Key:   "gid",
			Value: featureGroupID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	FeatureGroupPUT(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Replacing feature group failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Replacing feature group failed: content type is invalid")
	}

	output := &featuregroup.Group{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Replacing feature group failed: %s", err)
	}

	if output.Name != input.Name {
		t.Errorf("Replacing feature group failed: feature group name %s and %s unequal", output.Name, input.Name)
	}
}

func TestFeatureGroupDELETE(t *testing.T) {
	featureGroupID := featureGroupIDs[len(featureGroupIDs)-1]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: featureGroupID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	FeatureGroupDELETE(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Deleting location failed: unexpcted response code %v", resp.StatusCode)
	}

	removeFeatureGroupID(featureGroupID)
}
