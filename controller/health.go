package controller

import (
	"net/http"

	"github.com/tarkov-database/rest-api/model/api"
	"github.com/tarkov-database/rest-api/view"

	"github.com/julienschmidt/httprouter"
)

// HealthGET handles a GET request on the health endpoint
func HealthGET(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	h, err := api.GetHealth()

	if err != nil || !h.OK {
		view.RenderJSON(h, http.StatusInternalServerError, w)
	} else {
		view.RenderJSON(h, http.StatusOK, w)
	}
}
