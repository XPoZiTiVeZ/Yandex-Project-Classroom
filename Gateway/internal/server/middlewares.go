package server

import (
	"Classroom/Gateway/internal/courses"
	he "Classroom/Gateway/internal/errors"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func HandlerWrapper[T any](handler func(w http.ResponseWriter, r *http.Request) (error, any)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body T
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			slog.Debug("unmarshal error", slog.Any("error", err))

			he.InternalError(w)
			return
		}

		slog.Debug("incoming", slog.Any("struct", body))

		ctx := r.Context()
		ctx = context.WithValue(ctx, "body", body)
		err, resp := handler(w, r.WithContext(ctx))

		if err != nil {
			return
		}

		slog.Debug("outcoming", slog.Any("struct", resp))

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			slog.Debug("marshal error", slog.Any("error", err))

			he.InternalError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
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
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid format of authorization header", http.StatusUnauthorized)
			return
		}

		claims := &AuthClaims{}
		token, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (any, error) {
			return s.Config.Common.AuthJWTSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid access token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (s *Server) IsTeacher(next http.HandlerFunc) http.HandlerFunc {
	return s.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*AuthClaims)
		if !ok {
			he.NotAuthenticated(w)
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

func (s *Server) IsMember(next http.HandlerFunc) http.HandlerFunc {
	return s.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*AuthClaims)
		if !ok {
			he.NotAuthenticated(w)
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

func (s *Server) IsSuperUser(next http.HandlerFunc) http.HandlerFunc {
	return s.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*AuthClaims)
		if !ok {
			he.NotAuthenticated(w)
			return
		}

		if !claims.IsSuperUser {
			he.NotSuperUser(w)
			return
		}

		next.ServeHTTP(w, r)
	}))
}
