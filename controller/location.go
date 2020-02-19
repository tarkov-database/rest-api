package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/model/location"
	"github.com/tarkov-database/rest-api/model/location/feature"
	"github.com/tarkov-database/rest-api/model/location/featuregroup"
	"github.com/tarkov-database/rest-api/view"

	"github.com/google/logger"
	"github.com/julienschmidt/httprouter"
)

var errInvalidLocationID = errors.New("invalid location id")

// LocationGET handles a GET request on a location entity endpoint
func LocationGET(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	loc, err := location.GetByID(ps.ByName("id"))
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(loc, http.StatusOK, w)
}

// LocationsGET handles a GET request on the location root endpoint
func LocationsGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var result *model.Result
	var err error

	opts := &location.Options{Sort: getSort("-_modified", r)}
	opts.Limit, opts.Offset = getLimitOffset(r)

Loop:
	for p, v := range r.URL.Query() {
		switch p {
		case "text":
			txt, err := url.QueryUnescape(v[0])
			if err != nil {
				s := &Status{}
				s.BadRequest(fmt.Sprintf("Query string error: %s", err.Error())).Render(w)
				return
			}

			result, err = location.GetByText(txt, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		case "available":
			available, err := strconv.ParseBool(v[0])
			if err != nil {
				s := &Status{}
				s.BadRequest(err.Error()).Render(w)
				return
			}

			result, err = location.GetByAvailability(available, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		}
	}

	if result == nil {
		result, err = location.GetAll(opts)
		if err != nil {
			handleError(err, w)
			return
		}
	}

	view.RenderJSON(result, http.StatusOK, w)
}

// LocationPOST handles a POST request on the location root endpoint
func LocationPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	loc := &location.Location{}

	if err := parseJSONBody(r.Body, loc); err != nil {
		s := &Status{}
		s.BadRequest(fmt.Sprintf("JSON parsing error: %s", err.Error())).Render(w)
		return
	}

	if err := loc.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(fmt.Sprintf("Validation error: %s", err.Error())).Render(w)
		return
	}

	if err := location.Create(loc); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Location %s created", loc.ID.Hex())

	view.RenderJSON(loc, http.StatusCreated, w)
}

// LocationPUT handles a PUT request on a location entity endpoint
func LocationPUT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	loc := &location.Location{}

	if err := parseJSONBody(r.Body, loc); err != nil {
		s := &Status{}
		s.BadRequest(fmt.Sprintf("JSON parsing error: %s", err.Error())).Render(w)
		return
	}

	if err := loc.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(fmt.Sprintf("Validation error: %s", err.Error())).Render(w)
		return
	}

	id := ps.ByName("id")

	if !loc.ID.IsZero() && loc.ID.Hex() != id {
		s := &Status{}
		s.UnprocessableEntity("ID mismatch").Render(w)
		return
	}

	if err := location.Replace(id, loc); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Location %s updated", loc.ID.Hex())

	view.RenderJSON(loc, http.StatusOK, w)
}

// LocationDELETE handles a DELETE request on a location entity endpoint
func LocationDELETE(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if err := location.Remove(id); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Location %s removed", id)

	w.WriteHeader(http.StatusNoContent)
}

// FeatureGET handles a GET request on a feature entity endpoint
func FeatureGET(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	ft, err := feature.GetByID(ps.ByName("fid"), ps.ByName("id"))
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(ft, http.StatusOK, w)
}

// FeaturesGET handles a GET request on the feature root endpoint
func FeaturesGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var result *model.Result
	var err error

	opts := &feature.Options{Sort: getSort("-_modified", r)}
	opts.Limit, opts.Offset = getLimitOffset(r)

	lID := ps.ByName("id")

Loop:
	for p, v := range r.URL.Query() {
		switch p {
		case "text":
			txt, err := url.QueryUnescape(v[0])
			if err != nil {
				s := &Status{}
				s.BadRequest(fmt.Sprintf("Query string error: %s", err.Error())).Render(w)
				return
			}

			result, err = feature.GetByText(txt, lID, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		case "group":
			grp, err := url.QueryUnescape(v[0])
			if err != nil {
				s := &Status{}
				s.BadRequest(fmt.Sprintf("Query string error: %s", err.Error())).Render(w)
				return
			}

			result, err = feature.GetByGroup(grp, lID, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		}
	}

	if result == nil {
		result, err = feature.GetAll(ps.ByName("id"), opts)
		if err != nil {
			handleError(err, w)
			return
		}
	}

	view.RenderJSON(result, http.StatusOK, w)
}

// FeaturePOST handles a POST request on the feature root endpoint
func FeaturePOST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	lID := ps.ByName("id")

	ft := &feature.Feature{}

	if err := parseJSONBody(r.Body, ft); err != nil {
		s := &Status{}
		s.BadRequest(fmt.Sprintf("JSON parsing error: %s", err.Error())).Render(w)
		return
	}

	if err := ft.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(fmt.Sprintf("Validation error: %s", err.Error())).Render(w)
		return
	}

	if !ft.Location.IsZero() && ft.Location.Hex() != lID {
		s := &Status{}
		s.UnprocessableEntity("Location ID mismatch").Render(w)
		return
	}

	loc, err := location.GetByID(lID)
	if err != nil {
		s := &Status{}
		s.UnprocessableEntity("Location don't exist").Render(w)
		return
	}

	if ft.Location.IsZero() {
		ft.Location = loc.ID
	}

	if err := feature.Create(ft); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Feature %s created", ft.ID.Hex())

	view.RenderJSON(ft, http.StatusCreated, w)
}

// FeaturePUT handles a PUT request on a feature entity endpoint
func FeaturePUT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	fID := ps.ByName("fid")
	lID := ps.ByName("id")

	ft := &feature.Feature{}

	if err := parseJSONBody(r.Body, ft); err != nil {
		s := &Status{}
		s.BadRequest(fmt.Sprintf("JSON parsing error: %s", err.Error())).Render(w)
		return
	}

	if err := ft.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(fmt.Sprintf("Validation error: %s", err.Error())).Render(w)
		return
	}

	if !ft.Location.IsZero() && ft.Location.Hex() != lID {
		s := &Status{}
		s.UnprocessableEntity("Location mismatch").Render(w)
		return
	}

	if !ft.ID.IsZero() && ft.ID.Hex() != fID {
		s := &Status{}
		s.UnprocessableEntity("ID mismatch").Render(w)
		return
	}

	loc, err := location.GetByID(lID)
	if err != nil {
		s := &Status{}
		s.UnprocessableEntity("Location don't exist").Render(w)
		return
	}

	if ft.Location.IsZero() {
		ft.Location = loc.ID
	}

	if err := feature.Replace(fID, ft); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Feature %s updated", ft.ID.Hex())

	view.RenderJSON(ft, http.StatusOK, w)
}

// FeatureDELETE handles a DELETE request on a feature entity endpoint
func FeatureDELETE(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if err := feature.Remove(id); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Feature %s removed", id)

	w.WriteHeader(http.StatusNoContent)
}

// FeatureGroupGET handles a GET request on a feature group entity endpoint
func FeatureGroupGET(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	ft, err := featuregroup.GetByID(ps.ByName("gid"), ps.ByName("id"))
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(ft, http.StatusOK, w)
}

// FeatureGroupsGET handles a GET request on the feature group root endpoint
func FeatureGroupsGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var result *model.Result
	var err error

	opts := &featuregroup.Options{Sort: getSort("-_modified", r)}
	opts.Limit, opts.Offset = getLimitOffset(r)

	lID := ps.ByName("id")

Loop:
	for p, v := range r.URL.Query() {
		switch p {
		case "text":
			txt, err := url.QueryUnescape(v[0])
			if err != nil {
				s := &Status{}
				s.BadRequest(fmt.Sprintf("Query string error: %s", err.Error())).Render(w)
				return
			}

			result, err = featuregroup.GetByText(txt, lID, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		case "tag":
			q, err := url.QueryUnescape(v[0])
			if err != nil {
				s := &Status{}
				s.BadRequest(fmt.Sprintf("Query string error: %s", err.Error())).Render(w)
				return
			}

			tags := strings.Split(q, ",")

			result, err = featuregroup.GetByTags(tags, lID, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		}
	}

	if result == nil {
		result, err = featuregroup.GetAll(lID, opts)
		if err != nil {
			handleError(err, w)
			return
		}
	}

	view.RenderJSON(result, http.StatusOK, w)
}

// FeatureGroupPOST handles a POST request on the featuregroup root endpoint
func FeatureGroupPOST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	lID := ps.ByName("id")

	fg := &featuregroup.Group{}

	if err := parseJSONBody(r.Body, fg); err != nil {
		s := &Status{}
		s.BadRequest(fmt.Sprintf("JSON parsing error: %s", err.Error())).Render(w)
		return
	}

	if err := fg.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(fmt.Sprintf("Validation error: %s", err.Error())).Render(w)
		return
	}

	if !fg.Location.IsZero() && fg.Location.Hex() != lID {
		s := &Status{}
		s.UnprocessableEntity("Location ID mismatch").Render(w)
		return
	}

	loc, err := location.GetByID(lID)
	if err != nil {
		s := &Status{}
		s.UnprocessableEntity("Location don't exist").Render(w)
		return
	}

	if fg.Location.IsZero() {
		fg.Location = loc.ID
	}

	if err := featuregroup.Create(fg); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Feature group %s created", fg.ID.Hex())

	view.RenderJSON(fg, http.StatusCreated, w)
}

// FeatureGroupPUT handles a PUT request on a feature group entity endpoint
func FeatureGroupPUT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		s := &Status{}
		s.UnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	fID := ps.ByName("gid")
	lID := ps.ByName("id")

	fg := &featuregroup.Group{}

	if err := parseJSONBody(r.Body, fg); err != nil {
		s := &Status{}
		s.BadRequest(fmt.Sprintf("JSON parsing error: %s", err.Error())).Render(w)
		return
	}

	if err := fg.Validate(); err != nil {
		s := &Status{}
		s.UnprocessableEntity(fmt.Sprintf("Validation error: %s", err.Error())).Render(w)
		return
	}

	if !fg.Location.IsZero() && fg.Location.Hex() != lID {
		s := &Status{}
		s.UnprocessableEntity("Location mismatch").Render(w)
		return
	}

	if !fg.ID.IsZero() && fg.ID.Hex() != fID {
		s := &Status{}
		s.UnprocessableEntity("ID mismatch").Render(w)
		return
	}

	loc, err := location.GetByID(lID)
	if err != nil {
		s := &Status{}
		s.UnprocessableEntity("Location don't exist").Render(w)
		return
	}

	if fg.Location.IsZero() {
		fg.Location = loc.ID
	}

	if err := featuregroup.Replace(fID, fg); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Feature group %s updated", fg.ID.Hex())

	view.RenderJSON(fg, http.StatusOK, w)
}

// FeatureGroupDELETE handles a DELETE request on a feature group entity endpoint
func FeatureGroupDELETE(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if err := featuregroup.Remove(id); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Feature group %s removed", id)

	w.WriteHeader(http.StatusNoContent)
}
