package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/tarkov-database/rest-api/model"
	"go.mongodb.org/mongo-driver/bson"
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

func getSort(def string, r *http.Request) bson.D {
	sort := make(bson.D, 0, 1)

	sortStr := def
	if s := r.URL.Query().Get("sort"); len(s) > 1 {
		sortStr = s
	}

	if strings.HasPrefix(sortStr, "-") {
		sort = append(sort, bson.E{Key: strings.TrimPrefix(sortStr, "-"), Value: -1})
	} else {
		sort = append(sort, bson.E{Key: sortStr, Value: 1})
	}

	return sort
}

var regexNonAlnumBlankPunct = regexp.MustCompile(`[^[:alnum:][:blank:][:punct:]]`)

func isAlnumBlankPunct(s string) bool {
	return !regexNonAlnumBlankPunct.MatchString(s)
}

var regexNotAllowedQueryChars = regexp.MustCompile(`[^[:alnum:][:blank:]!#%&'()*+,\-./:;?_~]`)

func isAllowedQueryChars(s string) bool {
	return !regexNotAllowedQueryChars.MatchString(s)
}

func handleError(err error, w http.ResponseWriter) {
	var res *Status

	switch err {
	case model.ErrNoResult:
		res = StatusNotFound("Resource(s) not found")
	case model.ErrInvalidKind:
		res = StatusNotFound("Kind is not valid")
	case model.ErrInvalidObjectID:
		res = StatusNotFound("Resource ID is not valid")
	case model.ErrInvalidInput:
		res = StatusUnprocessableEntity("Input is not valid")
	case model.ErrInternalError:
		res = StatusInternalServerError("Backend error")
	default:
		res = StatusInternalServerError("Internal error")
	}

	res.Render(w)
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

func parseObjIDs(query string) ([]string, error) {
	q, err := url.QueryUnescape(query)
	if err != nil {
		return nil, fmt.Errorf("Query string error: %s", err)
	}

	if len(q) < 24 {
		return nil, fmt.Errorf("ID is not valid")
	}

	ids := strings.Split(q, ",")
	if len(ids) > 100 {
		return nil, fmt.Errorf("ID limit exceeded")
	}

	return ids, nil
}
