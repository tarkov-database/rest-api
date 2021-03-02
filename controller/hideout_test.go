package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tarkov-database/rest-api/model/hideout/module"
	"github.com/tarkov-database/rest-api/model/hideout/production"
	"github.com/tarkov-database/rest-api/model/item"

	"github.com/julienschmidt/httprouter"
)

type moduleResult struct {
	Count int64           `json:"total"`
	Items []module.Module `json:"items"`
}

func TestModuleGET(t *testing.T) {
	moduleID := moduleIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: moduleID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	ModuleGET(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting module failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting module failed: content type is invalid")
	}

	output := &module.Module{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Getting module failed: %s", err)
	}

	if output.ID != moduleID {
		t.Error("Getting module failed: module ID invalid")
	}
}

func TestModulesGET(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/v2/hideout/module", nil)

	w := httptest.NewRecorder()

	ModulesGET(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting modules failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting modules failed: content type is invalid")
	}

	res := &moduleResult{}

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		t.Fatalf("Getting modules failed: %s", err)
	}

	if res.Count < 1 {
		t.Error("Getting modules failed: result count invalid")
	}
	if len(res.Items) == 0 {
		t.Fatal("Getting modules failed: result empty")
	}
}

func TestModulePOST(t *testing.T) {
	moduleID := createModuleID()

	buf := new(bytes.Buffer)

	input := &module.Module{
		ID:     moduleID,
		Name:   "test module",
		Stages: make([]module.Stage, 1),
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Creating module failed: %s", err)
	}

	req := httptest.NewRequest("POST", "http://example.com/v2/hideout/module", buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	w := httptest.NewRecorder()

	ModulePOST(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Creating module failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Creating module failed: content type is invalid")
	}

	output := &module.Module{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Creating module failed: %s", err)
	}

	if output.ID != input.ID {
		t.Errorf("Creating module failed: module ID %s and %s unequal", output.ID, input.ID)
	}
}

func TestModulePUT(t *testing.T) {
	moduleID := moduleIDs[0]

	buf := new(bytes.Buffer)

	input := &module.Module{
		Name:   "modified test module",
		Stages: make([]module.Stage, 1),
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Replacing module failed: %s", err)
	}

	req := httptest.NewRequest("PUT", "http://example.com/v2/hideout/module", buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: moduleID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	ModulePUT(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Replacing module failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Replacing module failed: content type is invalid")
	}

	output := &module.Module{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Replacing module failed: %s", err)
	}

	if output.Name != input.Name {
		t.Errorf("Replacing module failed: module e-mail %s and %s unequal", output.Name, input.Name)
	}
}

func TestModuleDELETE(t *testing.T) {
	moduleID := moduleIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: moduleID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	ModuleDELETE(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Deleting module failed: unexpcted response code %v", resp.StatusCode)
	}

	removeModuleID(moduleID)
}

type productionResult struct {
	Count int64                   `json:"total"`
	Items []production.Production `json:"items"`
}

func TestProductionGET(t *testing.T) {
	productionID := productionIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: productionID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	ProductionGET(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting production failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting production failed: content type is invalid")
	}

	output := &production.Production{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Getting production failed: %s", err)
	}

	if output.ID != productionID {
		t.Error("Getting production failed: production ID invalid")
	}
}

func TestProductionsGET(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/v2/hideout/production", nil)

	w := httptest.NewRecorder()

	ProductionsGET(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting productions failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting productions failed: content type is invalid")
	}

	res := &productionResult{}

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		t.Fatalf("Getting productions failed: %s", err)
	}

	if res.Count < 1 {
		t.Error("Getting productions failed: result count invalid")
	}
	if len(res.Items) == 0 {
		t.Fatal("Getting productions failed: result empty")
	}
}

func TestProductionPOST(t *testing.T) {
	productionID := createProductionID()

	buf := new(bytes.Buffer)

	input := &production.Production{
		ID:      productionID,
		Module:  moduleIDs[0],
		Outcome: []production.ItemRef{{ID: itemIDs[0], Count: 3, Kind: item.KindCommon}},
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Creating production failed: %s", err)
	}

	req := httptest.NewRequest("POST", "http://example.com/v2/hideout/production", buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	w := httptest.NewRecorder()

	ProductionPOST(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Creating production failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Creating production failed: content type is invalid")
	}

	output := &production.Production{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Creating production failed: %s", err)
	}

	if output.ID != input.ID {
		t.Errorf("Creating production failed: production ID %s and %s unequal", output.ID, input.ID)
	}
}

func TestProductionPUT(t *testing.T) {
	productionID := productionIDs[0]

	buf := new(bytes.Buffer)

	input := &production.Production{
		Module:  moduleIDs[0],
		Outcome: []production.ItemRef{{ID: itemIDs[0], Count: 1, Kind: item.KindCommon}},
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Replacing production failed: %s", err)
	}

	req := httptest.NewRequest("PUT", "http://example.com/v2/hideout/production", buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: productionID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	ProductionPUT(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Replacing production failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Replacing production failed: content type is invalid")
	}

	output := &production.Production{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Replacing production failed: %s", err)
	}

	if output.Outcome[0].Count != input.Outcome[0].Count {
		t.Errorf("Replacing production failed: production e-mail %v and %v unequal",
			output.Outcome[0].Count, input.Outcome[0].Count)
	}
}

func TestProductionDELETE(t *testing.T) {
	productionID := productionIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: productionID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	ProductionDELETE(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Deleting production failed: unexpcted response code %v", resp.StatusCode)
	}

	removeProductionID(productionID)
}
