package server

import (
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Некорректный запрос"`
}

func mergeMsgs(message string, msgs []string) string {
	for _, msg := range msgs {
		message += ": " + msg
	}
	return message
}

// Так как у нас REST api, нужно чтобы сервис отдавал json, а http.Error пишет просто строку, и эти методы сделаем более общими, с возможностью дописать какое то сообщение

func InternalError(w http.ResponseWriter, msgs ...string) error {
	message := mergeMsgs("Internal server error", msgs)
	return WriteJSON(
		w,
		ErrorResponse{Code: http.StatusInternalServerError, Message: message},
		http.StatusInternalServerError,
	)
}

func Unauthorized(w http.ResponseWriter, msgs ...string) error {
	message := mergeMsgs("Unauthorized", msgs)
	return WriteJSON(
		w,
		ErrorResponse{Code: http.StatusUnauthorized, Message: message},
		http.StatusInternalServerError,
	)
}

func Forbidden(w http.ResponseWriter, msgs ...string) error {
	message := mergeMsgs("Forbidden", msgs)
	return WriteJSON(
		w,
		ErrorResponse{Code: http.StatusForbidden, Message: message},
		http.StatusInternalServerError,
	)
}

func BadRequest(w http.ResponseWriter, msgs ...string) error {
	message := mergeMsgs("Bad request", msgs)
	return WriteJSON(
		w,
		ErrorResponse{Code: http.StatusBadRequest, Message: message},
		http.StatusBadRequest,
	)
}

func AlreadyExists(w http.ResponseWriter, msgs ...string) error {
	message := mergeMsgs("Conflict", msgs)
	return WriteJSON(
		w,
		ErrorResponse{Code: http.StatusConflict, Message: message},
		http.StatusConflict,
	)
}

func NotFound(w http.ResponseWriter, msgs ...string) error {
	message := mergeMsgs("Not found", msgs)
	return WriteJSON(
		w,
		ErrorResponse{Code: http.StatusNotFound, Message: message},
		http.StatusNotFound,
	)
}

func ServiceUnavailable(w http.ResponseWriter, msgs ...string) error {
	message := mergeMsgs("Service unavailable", msgs)
	return WriteJSON(
		w,
		ErrorResponse{Code: http.StatusServiceUnavailable, Message: message},
		http.StatusServiceUnavailable,
	)
}
