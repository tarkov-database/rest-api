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

// Render renders Status
func (s *Status) Render(w http.ResponseWriter) {
	res := model.NewResponse(s.Message, s.Code)
	view.RenderJSON(res, s.Code, w)
}

// OK fills Status with an HTTP 200 status and message
func (s *Status) OK(msg string) *Status {
	s.Code, s.Message = http.StatusOK, msg

	return s
}

// Created fills Status with an HTTP 201 status and message
func (s *Status) Created(msg string) *Status {
	s.Code, s.Message = http.StatusCreated, msg

	return s
}

// BadRequest fills Status with an HTTP 400 status and message
func (s *Status) BadRequest(msg string) *Status {
	s.Code, s.Message = http.StatusBadRequest, msg

	return s
}

// Unauthorized fills Status with an HTTP 401 status and message
func (s *Status) Unauthorized(msg string) *Status {
	s.Code, s.Message = http.StatusUnauthorized, msg

	return s
}

// Forbidden fills Status with an HTTP 403 status and message
func (s *Status) Forbidden(msg string) *Status {
	s.Code, s.Message = http.StatusForbidden, msg

	return s
}

// NotFound fills Status with an HTTP 404 status and message
func (s *Status) NotFound(msg string) *Status {
	s.Code, s.Message = http.StatusNotFound, msg

	return s
}

// UnsupportedMediaType fills Status with an HTTP 415 status and message
func (s *Status) UnsupportedMediaType(msg string) *Status {
	s.Code, s.Message = http.StatusUnsupportedMediaType, msg

	return s
}

// UnprocessableEntity fills Status with an HTTP 422 status and message
func (s *Status) UnprocessableEntity(msg string) *Status {
	s.Code, s.Message = http.StatusUnprocessableEntity, msg

	return s
}

// InternalServerError fills Status with an HTTP 500 status and message
func (s *Status) InternalServerError(msg string) *Status {
	s.Code, s.Message = http.StatusInternalServerError, msg

	return s
}

// StatusNotFoundHandler returns a HTTP 404 handler
func StatusNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		s := &Status{}
		s.NotFound("Endpoint not found").Render(w)
	})
}
