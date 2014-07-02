package api

import (
	"encoding/json"
	"net/http"
)

func Json(data interface{}, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}
