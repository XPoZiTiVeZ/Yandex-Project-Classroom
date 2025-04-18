package repo

import (
	"Classroom/Courses/internal/domain"
	"database/sql"
	"time"
)

type Course struct {
	ID          string       `db:"course_id"`
	TeacherID   string       `db:"teacher_id"`
	Title       string       `db:"title"`
	Description string       `db:"description"`
	Visibility  bool         `db:"visibility"`
	StartTime   sql.NullTime `db:"start_time"`
	EndTime     sql.NullTime `db:"end_time"`
	CreatedAt   time.Time    `db:"created_at"`
}

func (c Course) ToDomain() domain.Course {
	var startTime, endTime *time.Time
	if c.StartTime.Valid {
		startTime = &c.StartTime.Time
	}
	if c.EndTime.Valid {
		endTime = &c.EndTime.Time
	}
	return domain.Course{
		ID:          c.ID,
		TeacherID:   c.TeacherID,
		Title:       c.Title,
		Description: c.Description,
		Visibility:  c.Visibility,
		StartTime:   startTime,
		EndTime:     endTime,
		CreatedAt:   c.CreatedAt,
	}
}

type Enrollment struct {
	StudentID  string    `db:"student_id"`
	CourseID   string    `db:"course_id"`
	EnrolledAt time.Time `db:"enrolled_at"`
}

func (e Enrollment) ToDomain() domain.Enrollment {
	return domain.Enrollment{
		StudentID:  e.StudentID,
		CourseID:   e.CourseID,
		EnrolledAt: e.EnrolledAt,
	}
}

type Student struct {
	UserID    string `db:"user_id"`
	Email     string `db:"email"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

func (m Student) ToDomain() domain.Student {
	return domain.Student{
		UserID:    m.UserID,
		Email:     m.Email,
		FirstName: m.FirstName,
		LastName:  m.LastName,
	}
}
