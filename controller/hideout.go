package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tarkov-database/rest-api/model"
	"github.com/tarkov-database/rest-api/model/hideout/module"
	"github.com/tarkov-database/rest-api/model/hideout/production"
	"github.com/tarkov-database/rest-api/view"

	"github.com/google/logger"
	"github.com/julienschmidt/httprouter"
)

// ModuleGET handles a GET request on a module entity endpoint
func ModuleGET(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	mod, err := module.GetByID(ps.ByName("id"))
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(mod, http.StatusOK, w)
}

// ModulesGET handles a GET request on the module root endpoint
func ModulesGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var result *model.Result
	var err error

	opts := &module.Options{Sort: getSort("-_modified", r)}
	opts.Limit, opts.Offset = getLimitOffset(r)

Loop:
	for p, v := range r.URL.Query() {
		switch p {
		case "id":
			q, err := url.QueryUnescape(v[0])
			if err != nil {
				StatusBadRequest(fmt.Sprintf("Query string error: %s", err)).Render(w)
				return
			}

			if len(q) < 24 {
				StatusBadRequest("ID is not valid").Render(w)
				return
			}

			ids := strings.Split(q, ",")
			if len(ids) > 100 {
				StatusBadRequest("ID limit exceeded").Render(w)
				return
			}

			result, err = module.GetByIDs(ids, opts)
			if err != nil {
				var res *Status

				switch err {
				case model.ErrInvalidInput:
					res = StatusUnprocessableEntity("Query contains an invalid ID")
				case model.ErrInternalError:
					res = StatusInternalServerError("Network or database error")
				default:
					res = StatusInternalServerError("Internal error")
				}

				res.Render(w)

				return
			}

			break Loop
		case "text":
			txt, err := url.QueryUnescape(v[0])
			if err != nil {
				StatusBadRequest(fmt.Sprintf("Query string error: %s", err)).Render(w)
				return
			}

			if l := len(txt); l < 3 || l > 32 {
				StatusBadRequest("Query string has an invalid length").Render(w)
				return
			}

			if !isAlnumBlankPunct(txt) {
				StatusBadRequest("Query string contains invalid characters").Render(w)
				return
			}

			result, err = module.GetByText(txt, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		case "material":
			mat, err := url.QueryUnescape(v[0])
			if err != nil {
				StatusBadRequest(fmt.Sprintf("Query string error: %s", err)).Render(w)
				return
			}

			result, err = module.GetByMaterial(mat, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		}
	}

	if result == nil {
		result, err = module.GetAll(opts)
		if err != nil {
			handleError(err, w)
			return
		}
	}

	view.RenderJSON(result, http.StatusOK, w)
}

// ModulePOST handles a POST request on the module root endpoint
func ModulePOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	mod := &module.Module{}

	if err := parseJSONBody(r.Body, mod); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := mod.Validate(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	if err := module.Create(mod); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Module %s created", mod.ID.Hex())

	view.RenderJSON(mod, http.StatusCreated, w)
}

// ModulePUT handles a PUT request on a module entity endpoint
func ModulePUT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	mod := &module.Module{}

	if err := parseJSONBody(r.Body, mod); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := mod.Validate(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	id := ps.ByName("id")

	if !mod.ID.IsZero() && mod.ID.Hex() != id {
		StatusUnprocessableEntity("ID mismatch").Render(w)
		return
	}

	if err := module.Replace(id, mod); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Module %s updated", mod.ID.Hex())

	view.RenderJSON(mod, http.StatusOK, w)
}

// ModuleDELETE handles a DELETE request on a module entity endpoint
func ModuleDELETE(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if err := module.Remove(id); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Module %s removed", id)

	w.WriteHeader(http.StatusNoContent)
}

// ProductionGET handles a GET request on a production entity endpoint
func ProductionGET(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	prod, err := production.GetByID(ps.ByName("id"))
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(prod, http.StatusOK, w)
}

// ProductionsGET handles a GET request on the production root endpoint
func ProductionsGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var result *model.Result
	var err error

	opts := &production.Options{Sort: getSort("-_modified", r)}
	opts.Limit, opts.Offset = getLimitOffset(r)

Loop:
	for p, v := range r.URL.Query() {
		switch p {
		case "id":
			q, err := url.QueryUnescape(v[0])
			if err != nil {
				StatusBadRequest(fmt.Sprintf("Query string error: %s", err)).Render(w)
				return
			}

			if len(q) < 24 {
				StatusBadRequest("ID is not valid").Render(w)
				return
			}

			ids := strings.Split(q, ",")
			if len(ids) > 100 {
				StatusBadRequest("ID limit exceeded").Render(w)
				return
			}

			result, err = production.GetByIDs(ids, opts)
			if err != nil {
				var res *Status

				switch err {
				case model.ErrInvalidInput:
					res = StatusUnprocessableEntity("Query contains an invalid ID")
				case model.ErrInternalError:
					res = StatusInternalServerError("Network or database error")
				default:
					res = StatusInternalServerError("Internal error")
				}

				res.Render(w)

				return
			}

			break Loop
		case "module":
			mod, err := url.QueryUnescape(v[0])
			if err != nil {
				StatusBadRequest(fmt.Sprintf("Query string error: %s", err)).Render(w)
				return
			}

			result, err = production.GetByModule(mod, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		case "material":
			mat, err := url.QueryUnescape(v[0])
			if err != nil {
				StatusBadRequest(fmt.Sprintf("Query string error: %s", err)).Render(w)
				return
			}

			result, err = production.GetByMaterial(mat, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		case "outcome":
			out, err := url.QueryUnescape(v[0])
			if err != nil {
				StatusBadRequest(fmt.Sprintf("Query string error: %s", err)).Render(w)
				return
			}

			result, err = production.GetByOutcome(out, opts)
			if err != nil {
				handleError(err, w)
				return
			}

			break Loop
		}
	}

	if result == nil {
		result, err = production.GetAll(opts)
		if err != nil {
			handleError(err, w)
			return
		}
	}

	view.RenderJSON(result, http.StatusOK, w)
}

// ProductionPOST handles a POST request on the production root endpoint
func ProductionPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	prod := &production.Production{}

	if err := parseJSONBody(r.Body, prod); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := prod.Validate(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	if err := production.Create(prod); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Production %s created", prod.ID.Hex())

	view.RenderJSON(prod, http.StatusCreated, w)
}

// ProductionPUT handles a PUT request on a production entity endpoint
func ProductionPUT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isSupportedMediaType(r) {
		StatusUnsupportedMediaType("Wrong content type").Render(w)
		return
	}

	prod := &production.Production{}

	if err := parseJSONBody(r.Body, prod); err != nil {
		StatusBadRequest(fmt.Sprintf("JSON parsing error: %s", err)).Render(w)
		return
	}

	if err := prod.Validate(); err != nil {
		StatusUnprocessableEntity(fmt.Sprintf("Validation error: %s", err)).Render(w)
		return
	}

	id := ps.ByName("id")

	if !prod.ID.IsZero() && prod.ID.Hex() != id {
		StatusUnprocessableEntity("ID mismatch").Render(w)
		return
	}

	if err := production.Replace(id, prod); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Production %s updated", prod.ID.Hex())

	view.RenderJSON(prod, http.StatusOK, w)
}

// ProductionDELETE handles a DELETE request on a production entity endpoint
func ProductionDELETE(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if err := production.Remove(id); err != nil {
		handleError(err, w)
		return
	}

	logger.Infof("Production %s removed", id)

	w.WriteHeader(http.StatusNoContent)
}
