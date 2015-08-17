package controller

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
)

func Json(status int, data interface{}, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	WriteCorsHeaders(w, r)
	w.WriteHeader(status)
	debug.PrintStack()
	return json.NewEncoder(w).Encode(data)
}
