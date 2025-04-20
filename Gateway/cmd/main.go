package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	srv "Classroom/Gateway/internal/server"
	cfg "Classroom/Gateway/pkg/config"
	"Classroom/Gateway/pkg/logger"
)

func main() {
	// slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()
	ctx = logger.NewDevelopment(ctx, logger.LevelDebug, false)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	config := cfg.MustReadConfig()

	server, err := srv.NewServer(ctx, config)
	server.CtxStop = stop
	if err != nil {
		logger.Error(ctx, "Server ran into problem: ", slog.Any("error", err))
		stop()
	}

	logger.Info(ctx, "Server running", slog.Int("port", config.Host.Port))
	go server.Run(ctx)

	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			fmt.Println()
		}

		server.Server.Shutdown(ctx)

		time.Sleep(300 * time.Millisecond)
		logger.Info(ctx, "Gracefully stopped")
	}
}
