package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/tarkov-database/rest-api/model/item"

	"github.com/julienschmidt/httprouter"
)

type itemResult struct {
	Count int64       `json:"total"`
	Items []item.Item `json:"items"`
}

func TestItemIndexGET(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/v2/item", nil)

	w := httptest.NewRecorder()

	ItemIndexGET(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting item failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting item failed: content type is invalid")
	}

	output := &item.Index{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Getting item failed: %s", err)
	}

	if output.Total == 0 {
		t.Error("Getting item index failed: total count is invalid")
	}
	if output.Modified.IsZero() {
		t.Error("Getting item index failed: global modified date is invalid")
	}

	if output.Kinds[item.KindCommon].Count == 0 {
		t.Error("Getting item index failed: kind count is invalid")
	}
	if output.Kinds[item.KindCommon].Modified.IsZero() {
		t.Error("Getting item index failed: kind modified date is invalid")
	}

	keyword := "item"

	val := url.Values{}
	val.Add("search", keyword)

	uri := fmt.Sprintf("http://example.com/v2/item?%s", val.Encode())
	req = httptest.NewRequest("GET", uri, nil)

	w = httptest.NewRecorder()

	ItemIndexGET(w, req, httprouter.Params{})

	resp = w.Result()
	defer resp.Body.Close()

	res := &itemResult{}

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		t.Fatalf("Getting items failed: %s", err)
	}

	if res.Count < 1 {
		t.Error("Getting items failed: result count invalid")
	}
	if len(res.Items) == 0 {
		t.Fatal("Getting items failed: result empty")
	}
	if name := res.Items[0].Name; !strings.HasPrefix(name, keyword) {
		t.Error("Getting items failed: item name prefix invalid")
	}
}

func TestItemGET(t *testing.T) {
	itemID := itemIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "kind",
			Value: "common",
		},
		httprouter.Param{
			Key:   "id",
			Value: itemID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	ItemGET(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting item failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting item failed: content type is invalid")
	}

	output := &item.Item{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Getting item failed: %s", err)
	}

	if output.ID != itemID {
		t.Error("Getting item failed: item ID invalid")
	}
}

func TestItemsGET(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/v2/item", nil)

	params := httprouter.Params{
		httprouter.Param{
			Key:   "kind",
			Value: "common",
		},
	}

	w := httptest.NewRecorder()

	ItemsGET(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting items failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting items failed: content type is invalid")
	}

	res := &itemResult{}

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		t.Fatalf("Getting items failed: %s", err)
	}

	if res.Count < 1 {
		t.Error("Getting items failed: result count invalid")
	}
	if len(res.Items) == 0 {
		t.Fatal("Getting items failed: result empty")
	}
}

func TestItemPOST(t *testing.T) {
	itemID := createItemID()

	buf := new(bytes.Buffer)

	input := &item.Item{
		ID:          itemID,
		Name:        "test item",
		ShortName:   "test",
		Description: "test description",
		Price:       1000,
		Weight:      3.7,
		MaxStack:    1,
		Rarity:      "rare",
		Kind:        "common",
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Creating item failed: %s", err)
	}

	req := httptest.NewRequest("POST", "http://example.com/v2/item", buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	params := httprouter.Params{
		httprouter.Param{
			Key:   "kind",
			Value: input.Kind.String(),
		},
	}

	w := httptest.NewRecorder()

	ItemPOST(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Creating item failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Creating item failed: content type is invalid")
	}

	output := &item.Item{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Creating item failed: %s", err)
	}

	if output.ID != itemID {
		t.Errorf("Creating item failed: item ID %s and %s unequal", output.ID, itemID)
	}
}

func TestItemPUT(t *testing.T) {
	itemID := itemIDs[0]

	buf := new(bytes.Buffer)

	input := &item.Item{
		ID:          itemID,
		Name:        "change item name",
		ShortName:   "test",
		Description: "test description",
		Price:       1000,
		Weight:      3.7,
		MaxStack:    1,
		Rarity:      "rare",
		Kind:        "common",
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Replacing item failed: %s", err)
	}

	req := httptest.NewRequest("PUT", "http://example.com/v2/item", buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	params := httprouter.Params{
		httprouter.Param{
			Key:   "kind",
			Value: input.Kind.String(),
		},
		httprouter.Param{
			Key:   "id",
			Value: itemID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	ItemPUT(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Replacing item failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Replacing item failed: content type is invalid")
	}

	output := &item.Item{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Replacing item failed: %s", err)
	}

	if output.Name != input.Name {
		t.Errorf("Replacing item failed: item name %s and %s unequal", output.Name, input.Name)
	}
}

func TestItemDELETE(t *testing.T) {
	itemID := itemIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: itemID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	ItemDELETE(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Deleting item failed: unexpcted response code %v", resp.StatusCode)
	}

	removeItemID(itemID)
}
