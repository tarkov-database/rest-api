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

// StatusOK fills Status with an HTTP 200 status and message
func StatusOK(msg string) *Status {
	return &Status{
		Code:    http.StatusOK,
		Message: msg,
	}
}

// StatusCreated fills Status with an HTTP 201 status and message
func StatusCreated(msg string) *Status {
	return &Status{
		Code:    http.StatusCreated,
		Message: msg,
	}
}

// StatusBadRequest fills Status with an HTTP 400 status and message
func StatusBadRequest(msg string) *Status {
	return &Status{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

// StatusUnauthorized fills Status with an HTTP 401 status and message
func StatusUnauthorized(msg string) *Status {
	return &Status{
		Code:    http.StatusUnauthorized,
		Message: msg,
	}
}

// StatusForbidden fills Status with an HTTP 403 status and message
func StatusForbidden(msg string) *Status {
	return &Status{
		Code:    http.StatusForbidden,
		Message: msg,
	}
}

// StatusNotFound fills Status with an HTTP 404 status and message
func StatusNotFound(msg string) *Status {
	return &Status{
		Code:    http.StatusNotFound,
		Message: msg,
	}
}

// StatusUnsupportedMediaType fills Status with an HTTP 415 status and message
func StatusUnsupportedMediaType(msg string) *Status {
	return &Status{
		Code:    http.StatusUnsupportedMediaType,
		Message: msg,
	}
}

// StatusUnprocessableEntity fills Status with an HTTP 422 status and message
func StatusUnprocessableEntity(msg string) *Status {
	return &Status{
		Code:    http.StatusUnprocessableEntity,
		Message: msg,
	}
}

// StatusInternalServerError fills Status with an HTTP 500 status and message
func StatusInternalServerError(msg string) *Status {
	return &Status{
		Code:    http.StatusInternalServerError,
		Message: msg,
	}
}

// StatusNotFoundHandler returns a HTTP 404 handler
func StatusNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		StatusNotFound("Endpoint not found").Render(w)
	})
}
