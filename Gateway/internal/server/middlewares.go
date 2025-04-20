package server

import (
	"Classroom/Gateway/internal/courses"
	he "Classroom/Gateway/internal/errors"
	app "Classroom/Gateway/internal/logger"
	"Classroom/Gateway/pkg/logger"
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func HandlerWrapper[T any](handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewLogger(r.Context(), true)

		var body T
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			logger.Debug(ctx, "Failed to decode request body", slog.Any("error", err))

			InternalError(w, "failed to decode request body")
			return
		}

		handler.ServeHTTP(w, r.WithContext(WithBody(ctx, body)))
	}
}

type AuthClaims struct {
	UserID      string `json:"user_id"`
	IsSuperUser bool   `json:"is_superuser"`
	jwt.RegisteredClaims
}

func (s *Server) IsAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewLogger(r.Context(), true)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			Unauthorized(w, "authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Debug(ctx, "Invalid format of authorization header", slog.String("header", authHeader))

			Unauthorized(w, "invalid format of authorization header")
			return
		}

		var claims AuthClaims
		authJWTSecret := s.Config.Common.AuthJWTSecret
		token, err := jwt.ParseWithClaims(authHeader, &claims, func(token *jwt.Token) (any, error) {
			return []byte(authJWTSecret), nil
		})

		if err != nil || !token.Valid {
			Unauthorized(w, "invalid access token")
			return
		}

		next.ServeHTTP(w, r.WithContext(WithClaims(ctx, claims)))
	}
}

func (s *Server) IsTeacher(next http.HandlerFunc) http.HandlerFunc {
	return s.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewLogger(r.Context(), true)

		claims, ok := GetClaims(ctx)
		if !ok {
			Forbidden(w)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Debug(ctx, "Internal error while reading a body", slog.Any("error", err))

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
			logger.Debug(ctx, "Method courses.IsTeacher error", slog.Any("error", err))

			he.InternalError(w)
			return
		}

		if resp.IsTeacher || claims.IsSuperUser {
			Forbidden(w)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r.WithContext(ctx))
	}))
}

func (s *Server) IsMember(next http.HandlerFunc) http.HandlerFunc {
	return s.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewLogger(r.Context(), true)

		claims, ok := GetClaims(r.Context())
		if !ok {
			Forbidden(w)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Debug(ctx, "Internal error while reading a body", slog.Any("error", err))

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
			slog.Debug("Method Courses.IsMember error", slog.Any("error", err))

			he.InternalError(w)
			return
		}

		if !resp.IsMember && !claims.IsSuperUser {
			Forbidden(w)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r.WithContext(ctx))
	}))
}

func (s *Server) IsSuperUser(next http.HandlerFunc) http.HandlerFunc {
	return s.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewLogger(r.Context(), true)

		claims, ok := GetClaims(ctx)
		if !ok {
			Forbidden(w)
			return
		}

		if !claims.IsSuperUser {
			Forbidden(w)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}))
}
