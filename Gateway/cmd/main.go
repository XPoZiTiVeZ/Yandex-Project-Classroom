package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	app "Classroom/Gateway/internal/logger"
	srv "Classroom/Gateway/internal/server"
	cfg "Classroom/Gateway/pkg/config"
	"Classroom/Gateway/pkg/logger"

	_ "Classroom/Gateway/docs"

	"github.com/joho/godotenv"
)

// @title Classroom Gateway API
// @version 1.0
// @description Gateway service for Classroom microservices. Handles authentication, routing, and request aggregation.
// @host localhost
// @BasePath /api
// @schemes http
// @produce json
// @consumes json
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter JWT token in format: "Bearer {token}"

func main() {
	ctx := context.Background()
	ctx = app.NewLogger(ctx, false)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	config := cfg.MustReadConfig()
	//TODO: не забыть убрать, пока что пусть будет для дебага
	fmt.Printf("%+v\n", config)

	server, err := srv.NewServer(ctx, config)
	server.CtxStop = stop
	if err != nil {
		logger.Error(ctx, "Server ran into problem: ", slog.Any("error", err))
		stop()
	}

	logger.Info(ctx, "Server running", slog.Int("port", config.Host.Port))
	go server.Run(ctx)

	<-ctx.Done()
	server.Server.Shutdown(ctx)
	logger.Info(ctx, "Gracefully stopped")
}

func init() {
	godotenv.Load()
}
