package controller

import (
	"net/http"
	"net/url"
	"strconv"
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
			s.BadRequest(err.Error()).Render(w)
			return
		}

		l, o := getLimitOffset(r)
		opts := &item.Options{Limit: l, Offset: o}

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

	l, o := getLimitOffset(r)

	opts := &item.Options{
		Sort:   getSort("_modified", -1, r),
		Limit:  l,
		Offset: o,
	}

Loop:
	for p, v := range r.URL.Query() {
		switch p {
		case "id":
			q, err := url.QueryUnescape(v[0])
			if err != nil {
				s := &Status{}
				s.BadRequest(err.Error()).Render(w)
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
				s.BadRequest(err.Error()).Render(w)
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
		qs := make(map[string]interface{})

		switch kind {
		case item.KindArmor:
			err = getQueryType(r.URL, qs)
			err = getQueryArmorClass(r.URL, qs)
		case item.KindFirearm:
			err = getQueryType(r.URL, qs)
			err = getQueryClass(r.URL, qs)
			err = getQueryCaliber(r.URL, qs)
		case item.KindTacticalrig:
			err = getQueryArmorClass(r.URL, qs)
		case item.KindAmmunition:
			err = getQueryType(r.URL, qs)
			err = getQueryCaliber(r.URL, qs)
		case item.KindMagazine:
			err = getQueryCaliber(r.URL, qs)
		case item.KindMedical, item.KindFood, item.KindGrenade, item.KindClothing, item.KindModificationMuzzle, item.KindModificationDevice, item.KindModificationSight, item.KindModificationSightSpecial, item.KindModificationGoggles, item.KindModificationGogglesSpecial:
			err = getQueryType(r.URL, qs)
		}
		if err != nil {
			s := &Status{}
			s.BadRequest(err.Error()).Render(w)
			return
		}

		result, err = item.GetAll(qs, kind, opts)
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
		s.BadRequest(err.Error()).Render(w)
		return
	}

	if err := entity.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(err.Error()).Render(w)
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
		s.BadRequest(err.Error()).Render(w)
		return
	}

	if err := entity.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(err.Error()).Render(w)
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

	err = item.Replace(id, entity)
	if err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Item %s updated", entity.GetID().Hex())

	view.RenderJSON(entity, http.StatusOK, w)
}

// ItemDELETE handles a DELETE request on a item entity endpoint
func ItemDELETE(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	err := item.Remove(id)
	if err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Item %s removed", id)

	w.WriteHeader(http.StatusNoContent)
}

func getQueryType(u *url.URL, qs map[string]interface{}) error {
	var err error
	if tp := u.Query().Get("type"); len(tp) > 0 {
		qs["type"], err = url.QueryUnescape(tp)
		if err != nil {
			return err
		}
	}

	return err
}

func getQueryCaliber(u *url.URL, qs map[string]interface{}) error {
	var err error
	if cal := u.Query().Get("caliber"); len(cal) > 0 {
		qs["caliber"], err = url.QueryUnescape(cal)
		if err != nil {
			return err
		}
	}

	return err
}

func getQueryClass(u *url.URL, qs map[string]interface{}) error {
	var err error
	if cl := u.Query().Get("class"); len(cl) > 0 {
		qs["class"], err = url.QueryUnescape(cl)
		if err != nil {
			return err
		}
	}

	return err
}

func getQueryArmorClass(u *url.URL, qs map[string]interface{}) error {
	if cl := u.Query().Get("class"); len(cl) > 0 {
		if cl, err := url.QueryUnescape(cl); err == nil {
			c, err := strconv.ParseInt(cl, 10, 64)
			if err != nil {
				return err
			}
			qs["armor"] = map[string]int64{"class": c}
		}
	}

	return nil
}
