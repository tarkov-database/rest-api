package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/model/item"
	"github.com/tarkov-database/rest-api/view"

	"github.com/google/logger"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
)

// ItemIndexGET handles a GET request on the item root endpoint
func ItemIndexGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var skipKinds bool

	if skip := r.URL.Query().Get("skipKinds"); len(skip) > 0 {
		if skip == "1" || skip == "true" {
			skipKinds = true
		}
	}

	idx, err := item.GetIndex(skipKinds)
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(idx, http.StatusOK, w)
}

// ItemGET handles a GET request on a item entity endpoint
func ItemGET(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	kind := item.Kind(ps.ByName("kind"))
	if !kind.IsValid() {
		StatusNotFound("Kind not found").Render(w)
		return
	}

	i, err := item.GetByID(ps.ByName("id"), kind)
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(i, http.StatusOK, w)
}

// ItemsGET handles a GET request on a item kind endpoint
func ItemsGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var result *model.Result
	var err error

	kind := item.Kind(ps.ByName("kind"))
	if !kind.IsValid() {
		StatusNotFound("Kind not found").Render(w)
		return
	}

	opts := &item.Options{Sort: getSort("-_modified", r)}
	opts.Limit, opts.Offset = getLimitOffset(r)

Loop:
	for p, v := range r.URL.Query() {
		switch p {
		case "id":
			q, err := url.QueryUnescape(v[0])
			if err != nil {
				StatusBadRequest(fmt.Sprintf("Query string error: %s", err)).Render(w)
				return
			}

			if len(q) < 24 {
				StatusBadRequest("ID is not valid").Render(w)
				return
			}

			ids := strings.Split(q, ",")
			if len(ids) > 100 {
				StatusBadRequest("ID limit exceeded").Render(w)
				return
			}

			result, err = item.GetByIDs(ids, kind, opts)
			if err != nil {
				var res *Status

				switch err {
				case model.ErrInvalidInput:
					res = StatusUnprocessableEntity("Query contains an invalid ID")
				case model.ErrInternalError:
					res = StatusInternalServerError("Network or database error")
				default:
					res = StatusInternalServerError("Internal error")
				}

				res.Render(w)

				return
			}

			break Loop
		case "text":
			txt, err := url.QueryUnescape(v[0])
			if err != nil {
				StatusBadRequest(fmt.Sprintf("Query string error: %s", err)).Render(w)
				return
			}

			if l := len(txt); l < 3 || l > 32 {
				StatusBadRequest("Query string has an invalid length").Render(w)
				return
			}

			if !isAlnumBlankPunct(txt) {
				StatusBadRequest("Query string contains invalid characters").Render(w)
				return
			}

			result, err = item.GetByText(txt, opts, kind)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		}
	}

	if result == nil {
		var filter model.DocumentFilter

		switch kind {
		case item.KindArmor:
			armorFilter := &item.ArmorFilter{}

			if v := r.URL.Query().Get("type"); v != "" {
				if !isAllowedQueryChars(v) {
					break
				}
				armorFilter.Type = &v
			}

			if v := r.URL.Query().Get("armor.class"); v != "" {
				if v, err := strconv.ParseInt(v, 10, 64); err == nil {
					armorFilter.ArmorClass = &v
				} else {
					break
				}
			}

			if v := r.URL.Query().Get("armor.material.name"); v != "" {
				if !isAllowedQueryChars(v) {
					break
				}
				armorFilter.MaterialName = &v
			}

			filter = armorFilter
		case item.KindFirearm:
			firearmFilter := &item.FirearmFilter{}

			if v := r.URL.Query().Get("type"); v != "" {
				if !isAllowedQueryChars(v) {
					break
				}
				firearmFilter.Type = &v
			}

			if v := r.URL.Query().Get("class"); v != "" {
				if !isAllowedQueryChars(v) {
					break
				}
				firearmFilter.Class = &v
			}

			if v := r.URL.Query().Get("caliber"); v != "" {
				if !isAllowedQueryChars(v) {
					break
				}
				firearmFilter.Caliber = &v
			}

			if v := r.URL.Query().Get("manufacturer"); v != "" {
				if !isAllowedQueryChars(v) {
					break
				}
				firearmFilter.Manufacturer = &v
			}

			filter = firearmFilter
		case item.KindTacticalrig:
			tacticalrigFilter := &item.TacticalRigFilter{}

			if v := r.URL.Query().Get("isPlateCarrier"); v != "" {
				if v, err := strconv.ParseBool(v); err == nil {
					tacticalrigFilter.IsPlateCarrier = &v
				} else {
					break
				}
			}

			if v := r.URL.Query().Get("isArmored"); v != "" {
				if v, err := strconv.ParseBool(v); err == nil {
					tacticalrigFilter.IsArmored = &v
				} else {
					break
				}
			}

			if v := r.URL.Query().Get("armor.class"); v != "" {
				if v, err := strconv.ParseInt(v, 10, 64); err == nil {
					tacticalrigFilter.ArmorClass = &v
				} else {
					break
				}
			}

			if v := r.URL.Query().Get("armor.material.name"); v != "" {
				if !isAllowedQueryChars(v) {
					break
				}
				tacticalrigFilter.ArmorMaterial = &v
			}

			filter = tacticalrigFilter
		case item.KindAmmunition:
			ammunitionFilter := &item.AmmunitionFilter{}

			if v := r.URL.Query().Get("type"); v != "" {
				if !isAllowedQueryChars(v) {
					break
				}
				ammunitionFilter.Type = &v
			}

			if v := r.URL.Query().Get("caliber"); v != "" {
				if !isAllowedQueryChars(v) {
					break
				}
				ammunitionFilter.Caliber = &v
			}

			filter = ammunitionFilter
		case item.KindMagazine:
			magazineFilter := &item.MagazineFilter{}

			if v := r.URL.Query().Get("caliber"); v != "" {
				if !isAllowedQueryChars(v) {
					break
				}
				magazineFilter.Caliber = &v
			}

			filter = magazineFilter
		case item.KindMedical, item.KindFood, item.KindGrenade, item.KindClothing, item.KindModificationMuzzle, item.KindModificationDevice, item.KindModificationSight, item.KindModificationSightSpecial, item.KindModificationGoggles:
			customFilter := bson.D{}

			if v := r.URL.Query().Get("type"); v != "" {
				if !isAllowedQueryChars(v) {
					break
				}
				customFilter = append(customFilter, bson.E{Key: "type", Value: v})
			}

			filter = &model.CustomFilter{D: customFilter}
		}
		if err != nil {
			StatusBadRequest(fmt.Sprintf("Query string error: %s", err)).Render(w)
			return
		}

		result, err = item.GetAll(filter, kind, opts)
		if err != nil {
			handleError(err, w)
			return
		}
	}

	view.RenderJSON(result, http.StatusOK, w)
}

// ItemPOST handles a POST request on a item kind endpoint
func ItemPOST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	kind := item.Kind(ps.ByName("kind"))

	entity, err := kind.GetEntity()
	if err != nil {
		handleError(err, w)
		return
	}

	if err := parseJSONBody(r.Body, entity); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := entity.Validate(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	if entity.GetKind() != kind {
		StatusUnprocessableEntity("Kind mismatch").Render(w)
		return
	}

	if err = item.Create(entity); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Item %s created", entity.GetID().Hex())

	view.RenderJSON(entity, http.StatusCreated, w)
}

// ItemPUT handles a PUT request on a item entity endpoint
func ItemPUT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	id, kind := ps.ByName("id"), item.Kind(ps.ByName("kind"))

	entity, err := kind.GetEntity()
	if err != nil {
		handleError(err, w)
		return
	}

	if err := parseJSONBody(r.Body, entity); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := entity.Validate(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	if docID := entity.GetID(); !docID.IsZero() && docID.Hex() != id {
		StatusUnprocessableEntity("ID mismatch").Render(w)
		return
	}

	if entity.GetKind() != kind {
		StatusUnprocessableEntity("Kind mismatch").Render(w)
		return
	}

	if err := item.Replace(id, entity); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Item %s updated", entity.GetID().Hex())

	view.RenderJSON(entity, http.StatusOK, w)
}

// ItemDELETE handles a DELETE request on a item entity endpoint
func ItemDELETE(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if err := item.Remove(id); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Item %s removed", id)

	w.WriteHeader(http.StatusNoContent)
}
