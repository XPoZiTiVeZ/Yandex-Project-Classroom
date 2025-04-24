package server

import (
	"Classroom/Gateway/internal/courses"
	app "Classroom/Gateway/internal/logger"
	"Classroom/Gateway/pkg/logger"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/schema"
)

func JSONHandlerWrapper[T any](handler http.HandlerFunc) http.HandlerFunc {
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

func QueryHandlerWrapper[T any](handler http.HandlerFunc)  http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewLogger(r.Context(), true)

		var body T
		err := schema.NewDecoder().Decode(&body, r.URL.Query())
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
		token, err := jwt.ParseWithClaims(parts[1], &claims, func(token *jwt.Token) (any, error) {
			return []byte(authJWTSecret), nil
		})

		if err != nil || !token.Valid {
			Unauthorized(w, "invalid access token")
			return
		}

		next.ServeHTTP(w, r.WithContext(WithClaims(ctx, claims)))
	}
}

func (s *Server) IsStudent(ctx context.Context, courseID string) (bool, error) {
	claims, _ := GetClaims(ctx)

	req := courses.IsMemberRequest{
		UserID: claims.UserID,
		CourseID: courseID,
	}

	resp, err := s.Courses.IsMember(ctx, s.Redis, req)
	if err != nil {
		return false, err
	}

	return resp.IsMember, nil
}

func (s *Server) IsTeacher(ctx context.Context, courseID string) (bool, error) {
	claims, _ := GetClaims(ctx)

	req := courses.IsTeacherRequest{
		UserID: claims.UserID,
		CourseID: courseID,
	}

	resp, err := s.Courses.IsTeacher(ctx, s.Redis, req)
	if err != nil {
		return false, err
	}

	return resp.IsTeacher, nil
}

func (s *Server) IsMember(ctx context.Context, courseID string) (bool, error) {
	claims, _ := GetClaims(ctx)

	req1 := courses.IsMemberRequest{
		UserID: claims.UserID,
		CourseID: courseID,
	}

	resp1, err := s.Courses.IsMember(ctx, s.Redis, req1)
	if err != nil {
		return false, err
	}

	req2 := courses.IsTeacherRequest{
		UserID: claims.UserID,
		CourseID: courseID,
	}

	resp2, err := s.Courses.IsTeacher(ctx, s.Redis, req2)
	if err != nil {
		return false, err
	}

	return resp1.IsMember || resp2.IsTeacher, nil
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
