package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/tarkov-database/rest-api/model"
)

func getLimitOffset(r *http.Request) (int64, int64) {
	limit, offset := 20, 0
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		if l > 100 {
			l = 100
		}
		limit = l
	}
	if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil {
		offset = o
	}

	return int64(limit), int64(offset)
}

func getSort(def string, order int64, r *http.Request) map[string]int64 {
	sort := map[string]int64{def: order}

	if s := r.URL.Query().Get("sort"); len(s) > 1 {
		if strings.HasPrefix(s, "-") {
			sort = map[string]int64{strings.TrimPrefix(s, "-"): -1}
		} else {
			sort = map[string]int64{s: 1}
		}
	}

	return sort
}

func handleError(err error, w http.ResponseWriter) {
	s := &Status{}
	switch err {
	case model.ErrNoResult:
		s.NotFound("Resource(s) not found").Render(w)
	case model.ErrInvalidKind:
		s.NotFound("Kind is not valid").Render(w)
	case model.ErrInvalidObjectID:
		s.NotFound("Resource ID is not valid").Render(w)
	case model.ErrInvalidInput:
		s.UnprocessableEntity("Input is not valid").Render(w)
	case model.ErrInternalError:
		s.InternalServerError("Backend error").Render(w)
	default:
		s.InternalServerError("Internal error").Render(w)
	}
}

func isSupportedMediaType(r *http.Request) bool {
	if r.Header.Get("Content-Type") == "application/json" {
		return true
	}

	return false
}

func parseJSONBody(body io.ReadCloser, target interface{}) error {
	err := json.NewDecoder(body).Decode(target)
	body.Close()

	return err
}
