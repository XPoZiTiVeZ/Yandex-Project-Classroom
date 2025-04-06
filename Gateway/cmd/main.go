package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	// _ "Classroom/Gateway/internal/auth"
	// _ "Classroom/Gateway/internal/course"
	// _ "Classroom/Gateway/internal/lesson"
	srv "Classroom/Gateway/internal/server"
	cfg "Classroom/Gateway/pkg/config"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	config, err := cfg.ReadConfig()
	if err != nil {
		slog.Error(err.Error())
	}

	server, err := srv.NewServer(config.Server.Address, config.Server.Port, ctx)
	server.CtxStop = stop
	server.Config = &config
	if err != nil {
		fmt.Println(123)
		slog.Error(err.Error())
		stop()
	}
	slog.Info(fmt.Sprintf("Server running on %s:%d", config.Server.Address, config.Server.Port))
	go server.Run()

	<-ctx.Done()
	fmt.Printf("Gracefully stopped\n")
}
