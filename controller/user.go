package controller

import (
	"errors"
	"net/http"
	"strconv"

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
	var usr interface{}
	var err error

	l, o := getLimitOffset(r)
	opts := &user.Options{
		Sort:   getSort("_modified", -1, r),
		Limit:  l,
		Offset: o,
	}

Loop:
	for p, v := range r.URL.Query() {
		switch p {
		case "locked":
			locked, err := strconv.ParseBool(v[0])
			if err != nil {
				s := &Status{}
				s.BadRequest(err.Error()).Render(w)
				return
			}

			usr, err = user.GetByLockedState(locked, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		case "email":
			usr, err = user.GetByEmail(v[0], opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		}
	}

	if usr == nil {
		usr, err = user.GetAll(opts)
		if err != nil {
			handleError(err, w)
			return
		}
	}

	view.RenderJSON(usr, http.StatusOK, w)
}

// UserPOST handles a POST request on the user root endpoint
func UserPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	usr := &user.User{}
	err := parseJSONBody(r.Body, usr)
	if err != nil {
		s := &Status{}
		s.BadRequest(err.Error()).Render(w)
		return
	}

	if err := usr.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(err.Error()).Render(w)
		return
	}

	err = user.Create(usr)
	if err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("User %s created", usr.ID.Hex())

	view.RenderJSON(usr, http.StatusCreated, w)
}

// UserPUT handles a PUT request on a user entity endpoint
func UserPUT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	usr := &user.User{}
	err := parseJSONBody(r.Body, usr)
	if err != nil {
		s := &Status{}
		s.BadRequest(err.Error()).Render(w)
		return
	}

	if err := usr.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(err.Error()).Render(w)
		return
	}

	id := ps.ByName("id")

	if !usr.ID.IsZero() && usr.ID.Hex() != id {
		s := &Status{}
		s.UnprocessableEntity("ID mismatch").Render(w)
		return
	}

	err = user.Replace(id, usr)
	if err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("User %s updated", usr.ID.Hex())

	view.RenderJSON(usr, http.StatusOK, w)
}

// UserDELETE handles a DELETE request on a user entity endpoint
func UserDELETE(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	err := user.Remove(ps.ByName("id"))
	if err != nil {
		handleError(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
