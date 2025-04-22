package service

import (
	"Classroom/Notifications/pkg/mailer"
	"context"
)

type UserRepo interface{}

type TaskRepo interface{}

type LessonRepo interface{}

type CourseRepo interface{}

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
	// TODO: отправить пользователю на почту что он был зачислен на курс
	return nil
}

func (s *notificationsService) UserExpelled(ctx context.Context, userID, courseID string) error {
	// TODO: отправить пользователю на почту что он был отчислен с курса
	return nil
}

func (s *notificationsService) LessonCreated(ctx context.Context, userID, courseID string) error {
	// TODO: оповестить всех пользователей с курса на котором был создан урок
	return nil
}

func (s *notificationsService) TaskCreated(ctx context.Context, userID, courseID string) error {
	// TODO: оповестить всех пользователей с курса на котором было создано задание
	return nil
}
