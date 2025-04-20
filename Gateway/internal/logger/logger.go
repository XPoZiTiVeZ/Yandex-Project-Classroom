package logger

import (
	"Classroom/Gateway/pkg/logger"
	"context"
)

func NewLogger(ctx context.Context, request bool) context.Context {
	return logger.NewDevelopment(ctx, logger.LevelDebug, request)
}
