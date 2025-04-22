package main

import (
	"Classroom/Notifications/internal/config"
	"Classroom/Notifications/internal/consumer"
	"Classroom/Notifications/pkg/events"
	"Classroom/Notifications/pkg/logger"
	"Classroom/Notifications/pkg/postgres"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx := logger.NewDevelopment(context.Background(), logger.LevelDebug, false)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	config := config.MustNew()
	fmt.Printf("%+v\n", config)

	consumer := consumer.MustNew([]string{config.KafkaBroker}, nil)
	defer consumer.Close()

	postgres := postgres.MustNew(config.PostgresURL)
	defer postgres.Close()

	consumer.ConsumeTopic(ctx, events.CourseEnrolledTopic)
	consumer.ConsumeTopic(ctx, events.CourseExpelledTopic)
	consumer.ConsumeTopic(ctx, events.TaskCreatedTopic)
	consumer.ConsumeTopic(ctx, events.LessonCreatedTopic)

	logger.Info(ctx, "Started")
	<-ctx.Done()
	logger.Info(ctx, "Stopped")
}

func init() {
	godotenv.Load()
}
