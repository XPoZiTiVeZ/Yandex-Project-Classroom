package server

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, resp any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(resp)
}
