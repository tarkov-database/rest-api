package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/model/user"
	"github.com/tarkov-database/rest-api/view"

	"github.com/google/logger"
	"github.com/julienschmidt/httprouter"
)

var errInvalidUserID = errors.New("invalid user id")

// UserGET handles a GET request on a user entity endpoint
func UserGET(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	usr, err := user.GetByID(ps.ByName("id"))
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(usr, http.StatusOK, w)
}

// UsersGET handles a GET request on the user root endpoint
func UsersGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var result *model.Result
	var err error

	opts := &user.Options{Sort: getSort("-_modified", r)}
	opts.Limit, opts.Offset = getLimitOffset(r)

Loop:
	for p, v := range r.URL.Query() {
		switch p {
		case "locked":
			s, err := url.QueryUnescape(v[0])
			if err != nil {
				StatusBadRequest(fmt.Sprintf("Query string error: %s", err)).Render(w)
				return
			}

			locked, err := strconv.ParseBool(s)
			if err != nil {
				StatusBadRequest(err.Error()).Render(w)
				return
			}

			result, err = user.GetByLockedState(locked, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		case "email":
			addr, err := url.QueryUnescape(v[0])
			if err != nil {
				StatusBadRequest(fmt.Sprintf("Query string error: %s", err)).Render(w)
				return
			}

			if l := len(addr); l < 3 || l > 100 {
				StatusBadRequest("Query string has an invalid length").Render(w)
				return
			}

			if !isAlnumBlankPunct(addr) {
				StatusBadRequest("Query string contains invalid characters").Render(w)
				return
			}

			result, err = user.GetByEmail(addr, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		}
	}

	if result == nil {
		result, err = user.GetAll(opts)
		if err != nil {
			handleError(err, w)
			return
		}
	}

	view.RenderJSON(result, http.StatusOK, w)
}

// UserPOST handles a POST request on the user root endpoint
func UserPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	usr := &user.User{}

	if err := parseJSONBody(r.Body, usr); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := usr.Validate(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	if err := user.Create(usr); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("User %s created", usr.ID.Hex())

	view.RenderJSON(usr, http.StatusCreated, w)
}

// UserPUT handles a PUT request on a user entity endpoint
func UserPUT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	usr := &user.User{}

	if err := parseJSONBody(r.Body, usr); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := usr.Validate(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	id := ps.ByName("id")

	if !usr.ID.IsZero() && usr.ID.Hex() != id {
		StatusUnprocessableEntity("ID mismatch").Render(w)
		return
	}

	if err := user.Replace(id, usr); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("User %s updated", usr.ID.Hex())

	view.RenderJSON(usr, http.StatusOK, w)
}

// UserDELETE handles a DELETE request on a user entity endpoint
func UserDELETE(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if err := user.Remove(id); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("User %s removed", id)

	w.WriteHeader(http.StatusNoContent)
}
