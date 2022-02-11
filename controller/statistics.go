package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/model/statistic/ammunition/armor"
	"github.com/tarkov-database/rest-api/model/statistic/ammunition/distance"
	"github.com/tarkov-database/rest-api/view"

	"github.com/google/logger"
	"github.com/julienschmidt/httprouter"
)

// DistanceStatGET handles a GET request on a ammunition statistics distance entity endpoint
func DistanceStatGET(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	loc, err := distance.GetByID(ps.ByName("id"))
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(loc, http.StatusOK, w)
}

// DistanceStatsGET handles a GET request on the distance root endpoint
func DistanceStatsGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var result *model.Result
	var err error

	opts := &distance.Options{Sort: getSort("distance", r)}
	opts.Limit, opts.Offset = getLimitOffset(r)

	var gte, lte *uint64
	if v := r.URL.Query().Get("range"); v != "" {
		v, err := url.QueryUnescape(v)
		if err != nil {
			StatusBadRequest("query value of \"range\" has an invalid value").Render(w)
			return
		}

		a := strings.SplitN(v, ",", 2)

		left, err := strconv.ParseUint(a[0], 10, 64)
		if err != nil {
			StatusBadRequest("query value of \"range\" has an invalid type").Render(w)
			return
		}

		gte = &left

		if len(a) > 1 {
			right, err := strconv.ParseUint(a[1], 10, 64)
			if err != nil {
				StatusBadRequest("query value of \"range\" has an invalid type").Render(w)
				return
			}

			lte = &right
		}
	}

	var ids []string
	if v := r.URL.Query().Get("ammo"); v != "" {
		v, err := parseObjIDs(v)
		if err != nil {
			StatusBadRequest(fmt.Sprintf("query value of \"ammo\" is invalid: %s", err)).Render(w)
			return
		}
		ids = v
	}

	result, err = distance.GetByRefsAndRange(ids, gte, lte, opts)
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(result, http.StatusOK, w)
}

// DistanceStatPOST handles a POST request on the distance root endpoint
func DistanceStatPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	stat := &distance.AmmoDistanceStatistics{}

	if err := parseJSONBody(r.Body, stat); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := stat.Validate(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	result, err := distance.GetByRefsAndRange([]string{stat.Reference.Hex()},
		&stat.Distance, &stat.Distance, &distance.Options{})
	if err != nil {
		handleError(err, w)
		return
	}

	if result.Count != 0 {
		StatusBadRequest("entity already exists").Render(w)
	}

	if err := distance.Create(stat); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Distance statistics %s created", stat.ID.Hex())

	view.RenderJSON(stat, http.StatusCreated, w)
}

// DistanceStatPUT handles a PUT request on a distance entity endpoint
func DistanceStatPUT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	stat := &distance.AmmoDistanceStatistics{}

	if err := parseJSONBody(r.Body, stat); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := stat.Validate(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	id := ps.ByName("id")

	if !stat.ID.IsZero() && stat.ID.Hex() != id {
		StatusUnprocessableEntity("ID mismatch").Render(w)
		return
	}

	if err := distance.Replace(id, stat); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Distance statistics %s updated", stat.ID.Hex())

	view.RenderJSON(stat, http.StatusOK, w)
}

// DistanceStatDELETE handles a DELETE request on a distance entity endpoint
func DistanceStatDELETE(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if err := distance.Remove(id); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Distance statistics %s removed", id)

	w.WriteHeader(http.StatusNoContent)
}

// ArmorStatGET handles a GET request on a ammunition statistics armor entity endpoint
func ArmorStatGET(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	loc, err := armor.GetByID(ps.ByName("id"))
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(loc, http.StatusOK, w)
}

// ArmorStatsGET handles a GET request on the distance root endpoint
func ArmorStatsGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var result *model.Result
	var err error

	opts := &armor.Options{Sort: getSort("distance", r)}
	opts.Limit, opts.Offset = getLimitOffset(r)

	rangeOpts := &armor.RangeOptions{}
	if v := r.URL.Query().Get("range"); v != "" {
		v, err := url.QueryUnescape(v)
		if err != nil {
			StatusBadRequest("query value of \"range\" has an invalid value").Render(w)
			return
		}

		a := strings.SplitN(v, ",", 2)

		gte, err := strconv.ParseUint(a[0], 10, 64)
		if err != nil {
			StatusBadRequest("query value of \"range\" has an invalid type").Render(w)
			return
		}

		rangeOpts.GTE = &gte

		if len(a) > 1 {
			lte, err := strconv.ParseUint(a[1], 10, 64)
			if err != nil {
				StatusBadRequest("query value of \"range\" has an invalid type").Render(w)
				return
			}

			rangeOpts.LTE = &lte
		}
	}

	var ammoIDs []string
	if v := r.URL.Query().Get("ammo"); v != "" {
		v, err := parseObjIDs(v)
		if err != nil {
			StatusBadRequest(fmt.Sprintf("query value of \"ammo\" is invalid: %s", err)).Render(w)
			return
		}
		ammoIDs = v
	}

	var armorIDs []string
	if v := r.URL.Query().Get("armor"); v != "" {
		v, err := parseObjIDs(v)
		if err != nil {
			StatusBadRequest(fmt.Sprintf("query value of \"armor\" is invalid: %s", err)).Render(w)
			return
		}
		armorIDs = v
	}

	result, err = armor.GetByRefs(ammoIDs, armorIDs, rangeOpts, opts)
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(result, http.StatusOK, w)
}

// ArmorStatPOST handles a POST request on the armor root endpoint
func ArmorStatPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	stat := &armor.AmmoArmorStatistics{}

	if err := parseJSONBody(r.Body, stat); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := stat.Validate(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	rangeOpts := &armor.RangeOptions{
		GTE: &stat.Distance,
		LTE: &stat.Distance,
	}

	result, err := armor.GetByRefs([]string{stat.Ammo.Hex()}, []string{stat.Armor.ID.Hex()},
		rangeOpts, &armor.Options{})
	if err != nil {
		handleError(err, w)
		return
	}

	if result.Count != 0 {
		StatusBadRequest("entity already exists").Render(w)
	}

	if err := armor.Create(stat); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Armor statistics %s created", stat.ID.Hex())

	view.RenderJSON(stat, http.StatusCreated, w)
}

// ArmorStatPUT handles a PUT request on a armor entity endpoint
func ArmorStatPUT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	stat := &armor.AmmoArmorStatistics{}

	if err := parseJSONBody(r.Body, stat); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := stat.Validate(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	id := ps.ByName("id")

	if !stat.ID.IsZero() && stat.ID.Hex() != id {
		StatusUnprocessableEntity("ID mismatch").Render(w)
		return
	}

	if err := armor.Replace(id, stat); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Armor statistics %s updated", stat.ID.Hex())

	view.RenderJSON(stat, http.StatusOK, w)
}

// ArmorStatDELETE handles a DELETE request on a armor entity endpoint
func ArmorStatDELETE(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if err := armor.Remove(id); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Armor statistics %s removed", id)

	w.WriteHeader(http.StatusNoContent)
}
