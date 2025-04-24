package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

type (
	loggerKey  string
	requestKey string
)

const (
	LoggerCtxKey loggerKey  = "logger"
	RequestIDKey requestKey = "request"
	LevelDebug              = slog.LevelDebug
	LevelError              = slog.LevelError
	LevelInfo               = slog.LevelInfo
)

func NewDevelopment(ctx context.Context, level slog.Level, request bool) context.Context {
	logger, ok := ctx.Value(LoggerCtxKey).(*slog.Logger)
	if ok {
		return ctx
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}
	logger = slog.New(slog.NewTextHandler(os.Stdout, opts))

	ctx = context.WithValue(ctx, LoggerCtxKey, logger)
	ctx = context.WithValue(ctx, RequestIDKey, request)

	return ctx
}

func NewProduction(ctx context.Context, request bool) context.Context {
	logger, ok := ctx.Value(LoggerCtxKey).(*slog.Logger)
	if ok {
		return ctx
	}

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))

	ctx = context.WithValue(ctx, LoggerCtxKey, logger)
	ctx = context.WithValue(ctx, RequestIDKey, request)

	return ctx
}

func getLoggerFromCtx(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(LoggerCtxKey).(*slog.Logger)
	if !ok {
		panic("No logger in context!")
	}

	return logger
}

func Info(ctx context.Context, msg string, fields ...any) {
	logger := getLoggerFromCtx(ctx)

	request := ctx.Value(RequestIDKey).(bool)
	requestID := "app"
	if request {
		requestID = uuid.New().String()
	}

	fields = append(fields, slog.String("flow", requestID))

	logger.Info(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...any) {
	logger := getLoggerFromCtx(ctx)

	request := ctx.Value(RequestIDKey).(bool)
	requestID := "app"
	if request {
		requestID = uuid.New().String()
	}

	fields = append(fields, slog.String("flow", requestID))

	logger.Info(msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...any) {
	logger := getLoggerFromCtx(ctx)

	request := ctx.Value(RequestIDKey).(bool)
	requestID := "app"
	if request {
		requestID = uuid.New().String()
	}

	fields = append(fields, slog.String("flow", requestID))

	logger.Info(msg, fields...)
}
