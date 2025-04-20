package http_errors

import "net/http"

func InternalError(w http.ResponseWriter) {
	http.Error(
		w,
		"Internal server error",
		http.StatusInternalServerError,
	)
}

func ServiceUnavailable(w http.ResponseWriter, msgs ...string) {
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

func NotSuperUser(w http.ResponseWriter) {
	http.Error(
		w,
		"Not a super user",
		http.StatusForbidden,
	)
}

func NotTacher(w http.ResponseWriter) {
	http.Error(
		w,
		"Not a teacher",
		http.StatusForbidden,
	)
}

func NotMember(w http.ResponseWriter) {
	http.Error(
		w,
		"Not a member",
		http.StatusForbidden,
	)
}

func NotAuthenticated(w http.ResponseWriter) {
	http.Error(
		w,
		"Not authenticated",
		http.StatusUnauthorized,
	)
}
