package server

import (
	"net/http"
	"strings"
)

// ErrorResponse стандартизированный формат для возврата ошибок API
// @description Cтандартизированный формат всех ошибок в API.
type ErrorResponse struct {
	// Код ошибки (HTTP статус код)
	Code int `json:"code" example:"400"`

	// Сообщение об ошибке
	Message string `json:"message" example:"Bad request"`
} // @name ErrorResponse
// @schema(title=ErrorResponse,required=["code","message"],order=["code","message"])

func mergeMsgs(message string, msgs []string) string {
	message += ": " + strings.Join(msgs, ": ")
	return message
}

func InternalError(w http.ResponseWriter, msgs ...string) {
	message := mergeMsgs("Internal server error", msgs)
	WriteJSON(
		w,
		ErrorResponse{Code: http.StatusInternalServerError, Message: message},
		http.StatusInternalServerError,
	)
}

func Unauthorized(w http.ResponseWriter, msgs ...string) {
	message := mergeMsgs("Unauthorized", msgs)
	WriteJSON(
		w,
		ErrorResponse{Code: http.StatusUnauthorized, Message: message},
		http.StatusInternalServerError,
	)
}

func Forbidden(w http.ResponseWriter, msgs ...string) {
	message := mergeMsgs("Forbidden", msgs)
	WriteJSON(
		w,
		ErrorResponse{Code: http.StatusForbidden, Message: message},
		http.StatusInternalServerError,
	)
}

func BadRequest(w http.ResponseWriter, msgs ...string) {
	message := mergeMsgs("Bad request", msgs)
	WriteJSON(
		w,
		ErrorResponse{Code: http.StatusBadRequest, Message: message},
		http.StatusBadRequest,
	)
}

func AlreadyExists(w http.ResponseWriter, msgs ...string) {
	message := mergeMsgs("Conflict", msgs)
	WriteJSON(
		w,
		ErrorResponse{Code: http.StatusConflict, Message: message},
		http.StatusConflict,
	)
}

func NotFound(w http.ResponseWriter, msgs ...string) {
	message := mergeMsgs("Not found", msgs)
	WriteJSON(
		w,
		ErrorResponse{Code: http.StatusNotFound, Message: message},
		http.StatusNotFound,
	)
}

func ServiceUnavailable(w http.ResponseWriter, msgs ...string) {
	message := mergeMsgs("Service unavailable", msgs)
	WriteJSON(
		w,
		ErrorResponse{Code: http.StatusServiceUnavailable, Message: message},
		http.StatusServiceUnavailable,
	)
}
