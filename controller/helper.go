package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/tarkov-database/rest-api/model"
)

func getLimitOffset(r *http.Request) (limit int64, offset int64) {
	limit = 20

	if l, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64); err == nil {
		if l > 100 {
			l = 100
		}
		limit = l
	}

	if o, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64); err == nil {
		offset = o
	}

	return
}

func getSort(def string, r *http.Request) map[string]int64 {
	sort := make(map[string]int64)

	sortStr := def
	if s := r.URL.Query().Get("sort"); len(s) > 1 {
		sortStr = s
	}

	if strings.HasPrefix(sortStr, "-") {
		sort = map[string]int64{strings.TrimPrefix(sortStr, "-"): -1}
	} else {
		sort = map[string]int64{sortStr: 1}
	}

	return sort
}

var regexNonAlnumBlankPunct = regexp.MustCompile(`[^[:alnum:][:blank:][:punct:]]`)

func isAlnumBlankPunct(s string) bool {
	return !regexNonAlnumBlankPunct.MatchString(s)
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
	defer body.Close()
	return json.NewDecoder(body).Decode(target)
}
