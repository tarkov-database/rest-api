package controller

import (
	"net/http"

	"github.com/tarkov-database/api/middleware/jwt"
	"github.com/tarkov-database/api/model/user"
	"github.com/tarkov-database/api/view"

	"github.com/julienschmidt/httprouter"
)

type Token struct {
	Token string `json:"token"`
}

func TokenGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token, err := jwt.GetToken(r)
	if err != nil {
		s := &Status{}
		s.Unauthorized(err.Error()).Write(w)
		return
	}

	clm, err := jwt.VerifyToken(token)
	if err != nil && err != jwt.ErrExpiredToken {
		s := &Status{}
		s.Unauthorized(err.Error()).Write(w)
		return
	}

	usr, err := user.GetByID(clm.Subject)
	if err != nil {
		handleError(err, w)
		return
	}

	if usr.Locked {
		s := &Status{}
		s.Forbidden("User is locked").Write(w)
		return
	}

	token, err = jwt.CreateToken(clm)
	if err != nil {
		s := &Status{}
		s.UnprocessableEntity(err.Error()).Write(w)
		return
	}

	view.RenderJSON(w, Token{token}, http.StatusCreated)
}

func TokenPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Write(w)
		return
	}

	issToken, err := jwt.GetToken(r)
	if err != nil {
		s := &Status{}
		s.Unauthorized(err.Error()).Write(w)
		return
	}

	issClaims, err := jwt.VerifyToken(issToken)
	if err != nil {
		s := &Status{}
		s.Unauthorized(err.Error()).Write(w)
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
		s.Forbidden("Insufficient permissions").Write(w)
		return
	}

	clm := &jwt.Claims{}
	err = parseJSONBody(r.Body, clm)
	if err != nil {
		s := &Status{}
		s.BadRequest(err.Error()).Write(w)
		return
	}

	if err := clm.ValidateCustom(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(err.Error()).Write(w)
		return
	}

	usr, err := user.GetByID(clm.Subject)
	if err != nil {
		handleError(err, w)
		return
	}

	if usr.Locked {
		s := &Status{}
		s.Forbidden("User is locked").Write(w)
		return
	}

	clm.Issuer = issClaims.Issuer

	token, err := jwt.CreateToken(clm)
	if err != nil {
		s := &Status{}
		s.InternalServerError(err.Error()).Write(w)
		return
	}

	view.RenderJSON(w, Token{token}, http.StatusCreated)
}
