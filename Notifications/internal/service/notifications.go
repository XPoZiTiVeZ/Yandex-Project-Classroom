package service

import (
	"Classroom/Notifications/internal/domain"
	"Classroom/Notifications/pkg/mailer"
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type UserRepo interface {
	GetByID(ctx context.Context, id string) (domain.User, error)
	ListByCourseID(ctx context.Context, courseID string) ([]domain.User, error)
}

type TaskRepo interface {
	GetByID(ctx context.Context, id string) (domain.Task, error)
}

type LessonRepo interface {
	GetByID(ctx context.Context, id string) (domain.Lesson, error)
}

type CourseRepo interface {
	GetByID(ctx context.Context, id string) (domain.Course, error)
}

type notificationsService struct {
	users   UserRepo
	tasks   TaskRepo
	lessons LessonRepo
	courses CourseRepo
	mailer  mailer.Mailer
}

func NewNotificationsService(mailer mailer.Mailer, users UserRepo, tasks TaskRepo, lessons LessonRepo, courses CourseRepo) *notificationsService {
	return &notificationsService{
		users:   users,
		tasks:   tasks,
		lessons: lessons,
		courses: courses,
		mailer:  mailer,
	}
}

func (s *notificationsService) UserEnrolled(ctx context.Context, userID, courseID string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	course, err := s.courses.GetByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to get course: %v", err)
	}

	subject := "Уведомление о зачислении на курс"
	body := fmt.Sprintf("Уважаемый %s %s, поздравляем с зачислением на курс %s, хорошего обучения.", user.FirstName, user.LastName, course.Title)
	return s.mailer.SendEmail(user.Email, subject, body)
}

func (s *notificationsService) UserExpelled(ctx context.Context, userID, courseID string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	course, err := s.courses.GetByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to get course: %v", err)
	}

	subject := "Уведомление об отчислении"
	body := fmt.Sprintf("Уважаемый %s %s, вы были отчислены с курса %s.", user.FirstName, user.LastName, course.Title)
	return s.mailer.SendEmail(user.Email, subject, body)
}

func (s *notificationsService) LessonCreated(ctx context.Context, lessonID, courseID string) error {
	users, err := s.users.ListByCourseID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}
	lesson, err := s.lessons.GetByID(ctx, lessonID)
	if err != nil {
		return fmt.Errorf("failed to get lesson: %v", err)
	}
	course, err := s.courses.GetByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to get course: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, user := range users {
		eg.Go(func() error {
			subject := fmt.Sprintf("На курс %s добавлен новый урок", course.Title)
			body := fmt.Sprintf(
				"%s %s, на курс %s, на котором вы обучаетесь был добавлен новый урок - %s",
				user.FirstName, user.LastName, course.Title, lesson.Title)

			return s.mailer.SendEmail(user.Email, subject, body)
		})
	}
	return eg.Wait()
}

func (s *notificationsService) TaskCreated(ctx context.Context, taskID, courseID string) error {
	users, err := s.users.ListByCourseID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}
	task, err := s.tasks.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %v", err)
	}
	course, err := s.courses.GetByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to get course: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, user := range users {
		eg.Go(func() error {
			subject := fmt.Sprintf("На курс %s добавлено новое задание", course.Title)
			body := fmt.Sprintf(
				"%s %s, на курс %s, на котором вы обучаетесь было добавлено новое задание - %s",
				user.FirstName, user.LastName, course.Title, task.Title)

			return s.mailer.SendEmail(user.Email, subject, body)
		})
	}
	return eg.Wait()
}
