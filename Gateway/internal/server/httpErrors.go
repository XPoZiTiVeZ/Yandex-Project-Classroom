package server

import "net/http"

func InternalError(w http.ResponseWriter) {
	http.Error(
		w,
		"Internal server error",
		http.StatusInternalServerError,
	)
}

func ServiceUnavailable(w http.ResponseWriter) {
	http.Error(
		w,
		"Service unavailable error",
		http.StatusServiceUnavailable,
	)
}

func AlreadyExists(w http.ResponseWriter) {
	http.Error(
		w,
		"User alredy exists",
		http.StatusConflict,
	)
}

func BadRequest(w http.ResponseWriter) {
	http.Error(
		w,
		"Invalid arguments",
		http.StatusBadRequest,
	)
}

func NotFound(w http.ResponseWriter) {
	http.Error(
		w,
		"Resource not found",
		http.StatusNotFound,
	)
}