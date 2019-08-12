package controller

import (
	"net/http"

	"github.com/tarkov-database/rest-api/model/api"
	"github.com/tarkov-database/rest-api/view"

	"github.com/julienschmidt/httprouter"
)

func IndexGET(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	i, err := api.GetIndex()
	if err != nil {
		handleError(err, w)
		return
	}

	view.RenderJSON(w, i, http.StatusOK)
}
