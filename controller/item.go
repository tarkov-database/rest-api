package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/model/item"
	"github.com/tarkov-database/rest-api/view"

	"github.com/google/logger"
	"github.com/julienschmidt/httprouter"
)

// ItemIndexGET handles a GET request on the item root endpoint
func ItemIndexGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var err error
	var i interface{}

	search := r.URL.Query().Get("search")
	switch {
	case len(search) > 0:
		txt, err := url.QueryUnescape(search)
		if err != nil {
			s := &Status{}
			s.BadRequest(fmt.Sprintf("Query string error: %s", err.Error())).Render(w)
			return
		}

		if !isAlnumBlankPunct(txt) {
			s := &Status{}
			s.BadRequest("Query string contains invalid characters").Render(w)
			return
		}

		opts := &item.Options{}
		opts.Limit, opts.Offset = getLimitOffset(r)

		i, err = item.GetByText(txt, opts)
		if err != nil {
			handleError(err, w)
			return
		}
	default:
		var skipKinds bool

		if skip := r.URL.Query().Get("skipKinds"); len(skip) > 0 {
			if skip == "1" {
				skipKinds = true
			}
		}

		i, err = item.GetIndex(skipKinds)
		if err != nil {
			handleError(err, w)
			return
		}
	}

	view.RenderJSON(i, http.StatusOK, w)
}

// ItemGET handles a GET request on a item entity endpoint
func ItemGET(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	kind := item.Kind(ps.ByName("kind"))
	if !kind.IsValid() {
		s := &Status{}
		s.NotFound("Kind not found").Render(w)
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
		s := &Status{}
		s.NotFound("Kind not found").Render(w)
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
				s := &Status{}
				s.BadRequest(fmt.Sprintf("Query string error: %s", err.Error())).Render(w)
				return
			}

			if len(q) < 24 {
				s := &Status{}
				s.BadRequest("ID is not valid").Render(w)
				return
			}

			ids := strings.Split(q, ",")
			if len(ids) > 100 {
				s := &Status{}
				s.BadRequest("ID limit exceeded").Render(w)
				return
			}

			result, err = item.GetByIDs(ids, kind, opts)
			if err != nil {
				s := &Status{}
				switch err {
				case model.ErrInvalidInput:
					s.UnprocessableEntity("Query contains an invalid ID").Render(w)
				case model.ErrInternalError:
					s.InternalServerError("Network or database error").Render(w)
				default:
					s.InternalServerError("Internal error").Render(w)
				}
				return
			}

			break Loop
		case "text":
			txt, err := url.QueryUnescape(v[0])
			if err != nil {
				s := &Status{}
				s.BadRequest(fmt.Sprintf("Query string error: %s", err.Error())).Render(w)
				return
			}

			if !isAlnumBlankPunct(txt) {
				s := &Status{}
				s.BadRequest("Query string contains invalid characters").Render(w)
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
		filter := model.Filter{}

		switch kind {
		case item.KindArmor:
			err = filter.AddString("type", r.URL.Query().Get("type"))
			err = filter.AddInt("armor.class", r.URL.Query().Get("armor.class"))
			err = filter.AddString("armor.material.name", r.URL.Query().Get("armor.material.name"))
		case item.KindFirearm:
			err = filter.AddString("type", r.URL.Query().Get("type"))
			err = filter.AddString("class", r.URL.Query().Get("class"))
			err = filter.AddString("caliber", r.URL.Query().Get("caliber"))
		case item.KindTacticalrig:
			err = filter.AddInt("armor.class", r.URL.Query().Get("armor.class"))
			err = filter.AddString("armor.material.name", r.URL.Query().Get("armor.material.name"))
		case item.KindAmmunition:
			err = filter.AddString("type", r.URL.Query().Get("type"))
			err = filter.AddString("caliber", r.URL.Query().Get("caliber"))
		case item.KindMagazine:
			err = filter.AddString("caliber", r.URL.Query().Get("caliber"))
		case item.KindMedical, item.KindFood, item.KindGrenade, item.KindClothing, item.KindModificationMuzzle, item.KindModificationDevice, item.KindModificationSight, item.KindModificationSightSpecial, item.KindModificationGoggles:
			err = filter.AddString("type", r.URL.Query().Get("type"))
		}
		if err != nil {
			s := &Status{}
			s.BadRequest(fmt.Sprintf("Query string error: %s", err.Error())).Render(w)
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
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	kind := item.Kind(ps.ByName("kind"))

	entity, err := kind.GetEntity()
	if err != nil {
		handleError(err, w)
		return
	}

	if err := parseJSONBody(r.Body, entity); err != nil {
		s := &Status{}
		s.BadRequest(fmt.Sprintf("JSON parsing error: %s", err.Error())).Render(w)
		return
	}

	if err := entity.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(fmt.Sprintf("Validation error: %s", err.Error())).Render(w)
		return
	}

	if entity.GetKind() != kind {
		s := &Status{}
		s.UnprocessableEntity("Kind mismatch").Render(w)
		return
	}

	err = item.Create(entity)
	if err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Item %s created", entity.GetID().Hex())

	view.RenderJSON(entity, http.StatusCreated, w)
}

// ItemPUT handles a PUT request on a item entity endpoint
func ItemPUT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	id, kind := ps.ByName("id"), item.Kind(ps.ByName("kind"))

	entity, err := kind.GetEntity()
	if err != nil {
		handleError(err, w)
		return
	}

	if err := parseJSONBody(r.Body, entity); err != nil {
		s := &Status{}
		s.BadRequest(fmt.Sprintf("JSON parsing error: %s", err.Error())).Render(w)
		return
	}

	if err := entity.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(fmt.Sprintf("Validation error: %s", err.Error())).Render(w)
		return
	}

	if docID := entity.GetID(); !docID.IsZero() && docID.Hex() != id {
		s := &Status{}
		s.UnprocessableEntity("ID mismatch").Render(w)
		return
	}

	if entity.GetKind() != kind {
		s := &Status{}
		s.UnprocessableEntity("Kind mismatch").Render(w)
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
