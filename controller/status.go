package controller

import (
	"net/http"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/view"
)

type Status struct {
	Code    int
	Message string
}

func (s *Status) New(status int, msg string) *Status {
	s.Code = status
	s.Message = msg

	return s
}

func (s *Status) Write(w http.ResponseWriter) {
	res := &model.Response{}
	res.New(s.Message, s.Code)
	view.RenderJSON(w, res, s.Code)
}

func (s *Status) NotFound(msg string) *Status {
	return s.New(http.StatusNotFound, msg)
}

func (s *Status) BadRequest(msg string) *Status {
	return s.New(http.StatusBadRequest, msg)
}

func (s *Status) UnsupportedMediaType(msg string) *Status {
	return s.New(http.StatusUnsupportedMediaType, msg)
}

func (s *Status) InternalServerError(msg string) *Status {
	return s.New(http.StatusInternalServerError, msg)
}

func (s *Status) UnprocessableEntity(msg string) *Status {
	return s.New(http.StatusUnprocessableEntity, msg)
}

func (s *Status) Unauthorized(msg string) *Status {
	return s.New(http.StatusUnauthorized, msg)
}

func (s *Status) Forbidden(msg string) *Status {
	return s.New(http.StatusForbidden, msg)
}

func (s *Status) Created(msg string) *Status {
	return s.New(http.StatusCreated, msg)
}

func (s *Status) OK(msg string) *Status {
	return s.New(http.StatusOK, msg)
}

func StatusNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		s := &Status{}
		s.NotFound("Endpoint not found").Write(w)
	})
}
