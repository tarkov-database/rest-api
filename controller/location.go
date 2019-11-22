package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/model/location"
	"github.com/tarkov-database/rest-api/view"

	"github.com/google/logger"
	"github.com/julienschmidt/httprouter"
)

var errInvalidLocationID = errors.New("invalid location id")

// LocationGET handles a GET request on a location entity endpoint
func LocationGET(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	loc, err := location.GetByID(ps.ByName("id"))
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(loc, http.StatusOK, w)
}

// LocationsGET handles a GET request on the location root endpoint
func LocationsGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var result *model.Result
	var err error

	opts := &location.Options{Sort: getSort("-_modified", r)}
	opts.Limit, opts.Offset = getLimitOffset(r)

Loop:
	for p, v := range r.URL.Query() {
		switch p {
		case "text":
			txt, err := url.QueryUnescape(v[0])
			if err != nil {
				s := &Status{}
				s.BadRequest(fmt.Sprintf("Query string error: %s", err.Error())).Render(w)
				return
			}

			result, err = location.GetByText(txt, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		case "available":
			available, err := strconv.ParseBool(v[0])
			if err != nil {
				s := &Status{}
				s.BadRequest(err.Error()).Render(w)
				return
			}

			result, err = location.GetByAvailability(available, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		}
	}

	if result == nil {
		result, err = location.GetAll(opts)
		if err != nil {
			handleError(err, w)
			return
		}
	}

	view.RenderJSON(result, http.StatusOK, w)
}

// LocationPOST handles a POST request on the location root endpoint
func LocationPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	loc := &location.Location{}

	if err := parseJSONBody(r.Body, loc); err != nil {
		s := &Status{}
		s.BadRequest(fmt.Sprintf("JSON parsing error: %s", err.Error())).Render(w)
		return
	}

	if err := loc.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(fmt.Sprintf("Validation error: %s", err.Error())).Render(w)
		return
	}

	if err := location.Create(loc); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Location %s created", loc.ID.Hex())

	view.RenderJSON(loc, http.StatusCreated, w)
}

// LocationPUT handles a PUT request on a location entity endpoint
func LocationPUT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	loc := &location.Location{}

	if err := parseJSONBody(r.Body, loc); err != nil {
		s := &Status{}
		s.BadRequest(fmt.Sprintf("JSON parsing error: %s", err.Error())).Render(w)
		return
	}

	if err := loc.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(fmt.Sprintf("Validation error: %s", err.Error())).Render(w)
		return
	}

	id := ps.ByName("id")

	if !loc.ID.IsZero() && loc.ID.Hex() != id {
		s := &Status{}
		s.UnprocessableEntity("ID mismatch").Render(w)
		return
	}

	if err := location.Replace(id, loc); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Location %s updated", loc.ID.Hex())

	view.RenderJSON(loc, http.StatusOK, w)
}

// LocationDELETE handles a DELETE request on a location entity endpoint
func LocationDELETE(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if err := location.Remove(id); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Location %s removed", id)

	w.WriteHeader(http.StatusNoContent)
}
