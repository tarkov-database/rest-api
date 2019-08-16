package view

import (
	"encoding/json"
	"net/http"

	"github.com/google/logger"
)

const contentTypeJSON = "application/json"

// RenderJSON encodes the input data into JSON and sends it as response
func RenderJSON(data interface{}, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		logger.Error(err)
	}
}
