package controller

import (
	"net/http"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/view"
)

// Status is used for status responses
type Status struct {
	Code    int
	Message string
}

// New fills the Status object
func (s *Status) New(status int, msg string) *Status {
	s.Code = status
	s.Message = msg

	return s
}

// Render renders Status
func (s *Status) Render(w http.ResponseWriter) {
	res := &model.Response{}
	res.New(s.Message, s.Code)
	view.RenderJSON(w, res, s.Code)
}

// NotFound fills Status with an HTTP 404 status and message
func (s *Status) NotFound(msg string) *Status {
	return s.New(http.StatusNotFound, msg)
}

// BadRequest fills Status with an HTTP 400 status and message
func (s *Status) BadRequest(msg string) *Status {
	return s.New(http.StatusBadRequest, msg)
}

// UnsupportedMediaType fills Status with an HTTP 415 status and message
func (s *Status) UnsupportedMediaType(msg string) *Status {
	return s.New(http.StatusUnsupportedMediaType, msg)
}

// InternalServerError fills Status with an HTTP 500 status and message
func (s *Status) InternalServerError(msg string) *Status {
	return s.New(http.StatusInternalServerError, msg)
}

// UnprocessableEntity fills Status with an HTTP 422 status and message
func (s *Status) UnprocessableEntity(msg string) *Status {
	return s.New(http.StatusUnprocessableEntity, msg)
}

// Unauthorized fills Status with an HTTP 401 status and message
func (s *Status) Unauthorized(msg string) *Status {
	return s.New(http.StatusUnauthorized, msg)
}

// Forbidden fills Status with an HTTP 403 status and message
func (s *Status) Forbidden(msg string) *Status {
	return s.New(http.StatusForbidden, msg)
}

// Created fills Status with an HTTP 201 status and message
func (s *Status) Created(msg string) *Status {
	return s.New(http.StatusCreated, msg)
}

// OK fills Status with an HTTP 200 status and message
func (s *Status) OK(msg string) *Status {
	return s.New(http.StatusOK, msg)
}

// StatusNotFoundHandler returns a HTTP 404 handler
func StatusNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		s := &Status{}
		s.NotFound("Endpoint not found").Render(w)
	})
}
