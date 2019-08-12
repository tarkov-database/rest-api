package view

import (
	"encoding/json"
	"net/http"

	"github.com/google/logger"
)

const contentTypeJSON = "application/json"

func RenderJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		logger.Error(err)
	}
}
