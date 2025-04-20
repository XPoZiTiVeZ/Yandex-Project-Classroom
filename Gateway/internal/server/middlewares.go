package server

import (
	"Classroom/Gateway/internal/courses"
	he "Classroom/Gateway/internal/errors"
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Обработку запросов и ответов будет удобнее делать в handler, иначе будет какой то непонятный статус 200
func HandlerWrapper[T any](handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body T
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			// TODO: add error log
			InternalError(w, "failed to decode request body")
			return
		}

		handler(w, r.WithContext(WithBody(r.Context(), body)))
	}
}

type AuthClaims struct {
	UserID      string `json:"user_id"`
	IsSuperUser bool   `json:"is_superuser"`
	jwt.RegisteredClaims
}

func (s *Server) IsAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			Unauthorized(w, "authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			Unauthorized(w, "invalid format of authorization header")
			return
		}

		var claims AuthClaims
		token, err := jwt.ParseWithClaims(authHeader, &claims, func(token *jwt.Token) (any, error) {
			// тут нужен обязательно массив байтов, иначе будет паника
			return []byte(s.Config.Common.AuthJWTSecret), nil
		})

		if err != nil || !token.Valid {
			Unauthorized(w, "invalid access token")
			return
		}

		next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), claims)))
	}
}

// в методах проверки роли, лучше не писать явную причину, достаточно просто forbidden
// TODO: пересмотреть
func (s *Server) IsTeacher(next http.HandlerFunc) http.HandlerFunc {
	return s.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetClaims(r.Context())
		if !ok {
			Forbidden(w)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			he.InternalError(w)
			return
		}

		var req courses.IsTeacherRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			he.BadRequest(w)
			return
		}
		req.UserID = claims.UserID

		resp, err := s.Courses.IsTeacher(r.Context(), req)
		if err != nil {
			slog.Debug("courses.IsTeacher error", slog.Any("error", err))

			he.InternalError(w)
			return
		}

		if resp.IsTeacher || claims.IsSuperUser {
			he.NotTacher(w)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	}))
}

// в методах проверки роли, лучше не писать явную причину, достаточно просто forbidden
// TODO: пересмотреть
func (s *Server) IsMember(next http.HandlerFunc) http.HandlerFunc {
	return s.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetClaims(r.Context())
		if !ok {
			Forbidden(w)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			he.InternalError(w)
			return
		}

		var req courses.IsMemberRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			he.BadRequest(w)
			return
		}
		req.UserID = claims.UserID

		resp, err := s.Courses.IsMember(r.Context(), req)
		if err != nil {
			slog.Debug("courses.IsMember error", slog.Any("error", err))

			he.InternalError(w)
			return
		}

		if !resp.IsMember && !claims.IsSuperUser {
			he.NotMember(w)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	}))
}

// в методах проверки роли, лучше не писать явную причину, достаточно просто forbidden
func (s *Server) IsSuperUser(next http.HandlerFunc) http.HandlerFunc {
	return s.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetClaims(r.Context())
		if !ok {
			Forbidden(w)
			return
		}

		if !claims.IsSuperUser {
			Forbidden(w)
			return
		}

		next.ServeHTTP(w, r)
	}))
}
