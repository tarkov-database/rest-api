package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/model/item"
	"github.com/tarkov-database/rest-api/model/statistic/ammunition/armor"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/julienschmidt/httprouter"
)

type ammoArmorStatsResult struct {
	Count int64                       `json:"total"`
	Items []armor.AmmoArmorStatistics `json:"items"`
}

func TestArmorStatsGET(t *testing.T) {
	statID := ammoArmorStatsIDs[0]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: statID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	ArmorStatGET(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting ammo armor statistics failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting ammo armor statistics failed: content type is invalid")
	}

	output := &armor.AmmoArmorStatistics{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Getting ammo armor statistics failed: %s", err)
	}

	if output.ID != statID {
		t.Error("Getting ammo armor statistics failed: ammo armor statistics ID invalid")
	}
}

func TestArmorStatssGET(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/v2/statistic/ammunition/armor", nil)

	w := httptest.NewRecorder()

	ArmorStatsGET(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Getting ammo armor statisticss failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Getting ammo armor statisticss failed: content type is invalid")
	}

	res := &ammoArmorStatsResult{}

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		t.Fatalf("Getting ammo armor statisticss failed: %s", err)
	}

	if res.Count < 1 {
		t.Error("Getting ammo armor statisticss failed: result count invalid")
	}
	if len(res.Items) == 0 {
		t.Fatal("Getting ammo armor statisticss failed: result empty")
	}
}

func TestArmorStatsPOST(t *testing.T) {
	statID := createStatisticAmmoArmorID()

	buf := new(bytes.Buffer)

	input := &armor.AmmoArmorStatistics{
		ID:   statID,
		Ammo: primitive.NewObjectID(),
		Armor: armor.ItemRef{
			ID:   primitive.NewObjectID(),
			Kind: item.KindArmor,
		},
		Distance:                  50,
		PenetrationChance:         [4]float64{},
		AverageShotsToDestruction: armor.Statistics{},
		AverageShotsTo50Damage:    armor.Statistics{},
		Modified:                  model.Timestamp{Time: time.Now()},
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Creating ammo armor statistics failed: %s", err)
	}

	req := httptest.NewRequest("POST", "http://example.com/v2/statistic/ammunition/armor", buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	w := httptest.NewRecorder()

	ArmorStatPOST(w, req, httprouter.Params{})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Creating ammo armor statistics failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Creating ammo armor statistics failed: content type is invalid")
	}

	output := &armor.AmmoArmorStatistics{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Creating ammo armor statistics failed: %s", err)
	}

	if output.ID != input.ID {
		t.Errorf("Creating ammo armor statistics failed: ammo armor statistics ID %s and %s unequal", output.ID, input.ID)
	}
}

func TestArmorStatsPUT(t *testing.T) {
	statID := ammoArmorStatsIDs[0]

	buf := new(bytes.Buffer)

	input := &armor.AmmoArmorStatistics{
		ID:   statID,
		Ammo: primitive.NewObjectID(),
		Armor: armor.ItemRef{
			ID:   primitive.NewObjectID(),
			Kind: item.KindArmor,
		},
		Distance:                  10,
		PenetrationChance:         [4]float64{},
		AverageShotsToDestruction: armor.Statistics{},
		AverageShotsTo50Damage:    armor.Statistics{},
		Modified:                  model.Timestamp{Time: time.Now()},
	}

	if err := json.NewEncoder(buf).Encode(input); err != nil {
		t.Fatalf("Replacing ammo armor statistics failed: %s", err)
	}

	req := httptest.NewRequest("PUT", "http://example.com/v2/statistic/ammunition/armor", buf)
	req.Header.Set("Content-Type", contentTypeJSON)

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: statID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	ArmorStatPUT(w, req, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Replacing ammo armor statistics failed: unexpcted response code %v", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentTypeJSON {
		t.Error("Replacing ammo armor statistics failed: content type is invalid")
	}

	output := &armor.AmmoArmorStatistics{}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		t.Fatalf("Replacing ammo armor statistics failed: %s", err)
	}

	if output.Distance != input.Distance {
		t.Errorf("Replacing ammo armor statistics failed: distance %v and %v unequal", output.Distance, input.Distance)
	}
}

func TestArmorStatsDELETE(t *testing.T) {
	statID := ammoArmorStatsIDs[len(ammoArmorStatsIDs)-1]

	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: statID.Hex(),
		},
	}

	w := httptest.NewRecorder()

	ArmorStatDELETE(w, &http.Request{}, params)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Deleting ammo armor statistics failed: unexpcted response code %v", resp.StatusCode)
	}

	removeStatisticAmmoArmorID(statID)
}
