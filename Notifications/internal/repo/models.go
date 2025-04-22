package repo

import "Classroom/Notifications/internal/domain"

type User struct {
	ID        string `db:"user_id"`
	Email     string `db:"email"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

func (u User) ToDomain() domain.User {
	return domain.User{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

type Course struct {
	ID    string `db:"course_id"`
	Title string `db:"title"`
}

func (c Course) ToDomain() domain.Course {
	return domain.Course{
		ID:    c.ID,
		Title: c.Title,
	}
}

type Task struct {
	ID    string `db:"task_id"`
	Title string `db:"title"`
}

func (t Task) ToDomain() domain.Task {
	return domain.Task{
		ID:    t.ID,
		Title: t.Title,
	}
}

type Lesson struct {
	ID    string `db:"lesson_id"`
	Title string `db:"title"`
}

func (l Lesson) ToDomain() domain.Lesson {
	return domain.Lesson{
		ID:    l.ID,
		Title: l.Title,
	}
}
