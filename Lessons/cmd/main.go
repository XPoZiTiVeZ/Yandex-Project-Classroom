package main

import (
	"Classroom/Lessons/internal/config"
	"Classroom/Lessons/pkg/postgres"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Для логгирования мы используем slog
	logger := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelDebug}))
	// Загрузку конфигурации лучше вынести в отдельный метод и пакет
	conf := config.MustNew()

	// Коннект к базе данных тоже лучше вынести в отдельный метод
	// Инициализацию таблиц не должно делать приложение, для этого используются миграции один раз, а не перед каждым запуском
	postgres := postgres.MustNew(conf.PostgresURL)
	defer postgres.Close()

	server := grpc.NewServer()

	// Для graceful shutdown лучше использовать контекст, так лаконичнее
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info("starting grpc server", "port", conf.Port)
	go startServer(server, conf.Port)

	<-ctx.Done()

	server.GracefulStop()
	logger.Info("grpc server stopped")
}

func startServer(server *grpc.Server, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// log Fatal завершает программу через os.Exit
		log.Fatalf("failed to listen: %v", err)
	}
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}
}

// Загружаем переменные окружения перед стартом
func init() {
	godotenv.Load()
}
