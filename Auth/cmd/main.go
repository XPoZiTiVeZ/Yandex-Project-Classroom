package main

import (
	"Classroom/Auth/internal/app"
	"Classroom/Auth/internal/config"
	"Classroom/Auth/internal/controller"
	"Classroom/Auth/internal/repo"
	"Classroom/Auth/internal/service"
	"Classroom/Auth/pkg/postgres"
	"Classroom/Auth/pkg/redis"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelDebug}))
	conf := config.MustNew()

	postgres := postgres.MustNew(conf.PostgresURL)
	defer postgres.Close()

	redis := redis.MustNew(conf.RedisURL)
	defer redis.Close()

	userRepo := repo.NewUserRepo(postgres)
	tokenRepo := repo.NewTokenRepo(redis)
	authService := service.NewAuthService(logger, userRepo, tokenRepo, conf.Auth)
	authController := controller.NewAuthController(logger, authService)

	app := app.New(logger, conf, authController)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// запускаем сервер в отдельной горутине и ждем завершения контекста
	go app.Start()
	<-ctx.Done()
	app.Stop()
}

// Загружаем переменные окружения перед стартом, не проверяем ошибку потому что не всегда будет .env файл
func init() {
	godotenv.Load()
}
