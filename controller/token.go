package controller

import (
	"fmt"
	"net/http"

	"github.com/tarkov-database/rest-api/middleware/jwt"
	"github.com/tarkov-database/rest-api/model/user"
	"github.com/tarkov-database/rest-api/view"

	"github.com/julienschmidt/httprouter"
)

// Token represents the body of a token creation response
type Token struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

// TokenGET handles a GET request on the token endpoint
func TokenGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token, err := jwt.GetToken(r)
	if err != nil {
		s := &Status{}
		s.Unauthorized(err.Error()).Render(w)
		return
	}

	clm, err := jwt.VerifyToken(token)
	if err != nil && err != jwt.ErrExpiredToken {
		s := &Status{}
		s.Unauthorized(err.Error()).Render(w)
		return
	}

	usr, err := user.GetByID(clm.Subject)
	if err != nil {
		handleError(err, w)
		return
	}

	if usr.Locked {
		s := &Status{}
		s.Forbidden("User is locked").Render(w)
		return
	}

	token, err = jwt.CreateToken(clm)
	if err != nil {
		s := &Status{}
		s.UnprocessableEntity(fmt.Sprintf("Creation error: %s", err)).Render(w)
		return
	}

	view.RenderJSON(Token{token, clm.ExpirationTime.Unix()}, http.StatusCreated, w)
}

// TokenPOST handles a POST request on the token endpoint
func TokenPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	issToken, err := jwt.GetToken(r)
	if err != nil {
		s := &Status{}
		s.Unauthorized(err.Error()).Render(w)
		return
	}

	issClaims, err := jwt.VerifyToken(issToken)
	if err != nil {
		s := &Status{}
		s.Unauthorized(err.Error()).Render(w)
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
		s := &Status{}
		s.Forbidden("Insufficient permissions").Render(w)
		return
	}

	clm := &jwt.Claims{}

	if err := parseJSONBody(r.Body, clm); err != nil {
		s := &Status{}
		s.BadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := clm.ValidateCustom(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	usr, err := user.GetByID(clm.Subject)
	if err != nil {
		handleError(err, w)
		return
	}

	if usr.Locked {
		s := &Status{}
		s.Forbidden("User is locked").Render(w)
		return
	}

	clm.Issuer = issClaims.Issuer

	token, err := jwt.CreateToken(clm)
	if err != nil {
		s := &Status{}
		s.InternalServerError(fmt.Sprintf("Creation error: %s", err)).Render(w)
		return
	}

	view.RenderJSON(Token{token, clm.ExpirationTime.Unix()}, http.StatusCreated, w)
}
