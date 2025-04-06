package server

import (
	"Classroom/Gateway/internal/auth"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong!"))
}

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var body auth.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Internal error while unmarshaling json: %v", err),
			http.StatusInternalServerError,
		)
	}

	slog.Debug("Incoming", slog.Any("struct", body))
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var body auth.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Internal error while unmarshaling json: %v", err),
			http.StatusInternalServerError,
		)
	}

	slog.Debug("Incoming", slog.Any("struct", body))
}
