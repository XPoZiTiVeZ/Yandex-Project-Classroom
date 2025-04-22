package main

import (
	"Classroom/Notifications/internal/config"
	"Classroom/Notifications/internal/consumer"
	"Classroom/Notifications/internal/repo"
	"Classroom/Notifications/internal/service"
	"Classroom/Notifications/pkg/events"
	"Classroom/Notifications/pkg/logger"
	"Classroom/Notifications/pkg/mailer"
	"Classroom/Notifications/pkg/postgres"
	"context"
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

	postgres := postgres.MustNew(config.PostgresURL)
	defer postgres.Close()

	mailer := mailer.New(mailer.Config{
		User: config.SMTP.User,
		Pass: config.SMTP.Pass,
		Port: config.SMTP.Port,
		Host: config.SMTP.Host,
	})
	userRepo := repo.NewUserRepo(postgres)
	courseRepo := repo.NewCourseRepo(postgres)
	taskRepo := repo.NewTaskRepo(postgres)
	lessonRepo := repo.NewLessonRepo(postgres)
	service := service.NewNotificationsService(mailer, userRepo, taskRepo, lessonRepo, courseRepo)

	consumer := consumer.MustNew([]string{config.KafkaBroker}, service)
	defer consumer.Close()

	consumer.ConsumeTopic(ctx, events.CourseEnrolledTopic)
	consumer.ConsumeTopic(ctx, events.CourseExpelledTopic)
	consumer.ConsumeTopic(ctx, events.TaskCreatedTopic)
	consumer.ConsumeTopic(ctx, events.LessonCreatedTopic)

	logger.Info(ctx, "service started")
	<-ctx.Done()
	logger.Info(ctx, "service stopped")
}

func init() {
	godotenv.Load()
}
