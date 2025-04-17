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
)

func main() {
	// slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	config := cfg.MustReadConfig()
	fmt.Println(*config)

	server, err := srv.NewServer(ctx, config)
	server.CtxStop = stop
	if err != nil {
		slog.Error("Server ran into problem: ", slog.Any("error", err))
		stop()
	}
	slog.Info(fmt.Sprintf("Server running on 0.0.0.0:%d", config.Host.Port))
	go server.Run()
	
	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			fmt.Println()
		}

		server.Server.Shutdown(ctx)

		time.Sleep(300 * time.Millisecond)
		slog.Info("Gracefully stopped\n")
	}
}
