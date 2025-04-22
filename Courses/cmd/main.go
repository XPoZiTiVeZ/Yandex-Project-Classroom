package main

import (
	"Classroom/Courses/internal/config"
	"Classroom/Courses/internal/producer"
	"Classroom/Courses/internal/repo"
	"Classroom/Courses/internal/service"
	pb "Classroom/Courses/pkg/api/courses"
	"Classroom/Courses/pkg/postgres"
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
	logger := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelDebug}))
	conf := config.MustNew()

	postgres := postgres.MustNew(conf.PostgresURL)
	defer postgres.Close()

	producer := producer.MustNewProducer([]string{conf.KafkaBroker})

	courseRepo := repo.NewCoursesRepo(postgres)
	courseService := service.NewCoursesService(logger, courseRepo, producer)

	server := grpc.NewServer()
	pb.RegisterCoursesServiceServer(server, courseService)

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
		log.Fatalf("failed to listen: %v", err)
	}
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}
}

func init() {
	godotenv.Load()
}
