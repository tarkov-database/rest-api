package controller

import (
	"fmt"
	"net/http"

	"github.com/tarkov-database/rest-api/middleware/jwt"
	"github.com/tarkov-database/rest-api/model/token"
	"github.com/tarkov-database/rest-api/model/user"
	"github.com/tarkov-database/rest-api/view"

	"github.com/julienschmidt/httprouter"
)

// TokenGET handles a GET request on the token endpoint
func TokenGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := jwt.GetToken(r)
	if err != nil {
		jwt.AddAuthenticateHeader(w, err)
		StatusUnauthorized(err.Error()).Render(w)
		return
	}

	clm, err := jwt.VerifyToken(t)
	if err != nil && err != jwt.ErrExpiredToken {
		jwt.AddAuthenticateHeader(w, err)
		StatusUnauthorized(err.Error()).Render(w)
		return
	}

	usr, err := user.GetByID(clm.Subject)
	if err != nil {
		handleError(err, w)
		return
	}

	if usr.Locked {
		StatusForbidden("User is locked").Render(w)
		return
	}

	t, err = jwt.CreateToken(clm, nil)
	if err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Creation error: %s", err)).Render(w)
		return
	}

	view.RenderJSON(token.Response{Token: t, Expires: clm.ExpirationTime.Unix()}, http.StatusCreated, w)
}

// TokenPOST handles a POST request on the token endpoint
func TokenPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	issToken, err := jwt.GetToken(r)
	if err != nil {
		jwt.AddAuthenticateHeader(w, err, jwt.ScopeTokenWrite, jwt.ScopeAllWrite)
		StatusUnauthorized(err.Error()).Render(w)
		return
	}

	issClaims, err := jwt.VerifyToken(issToken)
	if err != nil {
		jwt.AddAuthenticateHeader(w, err, jwt.ScopeTokenWrite, jwt.ScopeAllWrite)
		StatusUnauthorized(err.Error()).Render(w)
		return
	}

	var ok bool
	for _, s := range issClaims.Scope {
		if s == jwt.ScopeTokenWrite || s == jwt.ScopeAllWrite {
			ok = true
			break
		}
	}

	if !ok {
		jwt.AddAuthenticateHeader(w, jwt.ErrInvalidScope, jwt.ScopeTokenWrite, jwt.ScopeAllWrite)
		StatusForbidden("Insufficient permissions").Render(w)
		return
	}

	rb := &token.Request{}

	if err := parseJSONBody(r.Body, rb); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	lt, err := rb.Duration()
	if err != nil {
		StatusBadRequest(fmt.Sprintf("Parsing error: %s", err)).Render(w)
		return
	}

	clm := rb.ToClaims()

	if err := clm.ValidateCustom(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	usr, err := user.GetByID(clm.Subject)
	if err != nil {
		handleError(err, w)
		return
	}

	if usr.Locked {
		StatusForbidden("User is locked").Render(w)
		return
	}

	clm.Issuer = issClaims.Issuer

	t, err := jwt.CreateToken(clm, lt)
	if err != nil {
		StatusInternalServerError(fmt.Sprintf("Creation error: %s", err)).Render(w)
		return
	}

	view.RenderJSON(token.Response{Token: t, Expires: clm.ExpirationTime.Unix()}, http.StatusCreated, w)
}
