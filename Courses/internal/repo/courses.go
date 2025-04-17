package repo

import (
	"Classroom/Courses/internal/domain"
	"Classroom/Courses/internal/dto"
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type courseRepo struct {
	storage *sqlx.DB
	qb      sq.StatementBuilderType
}

func NewCoursesRepo(db *sqlx.DB) *courseRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &courseRepo{
		storage: db,
		qb:      qb,
	}
}

func (r *courseRepo) Create(ctx context.Context, dto dto.CreateCourseDTO) (domain.Course, error) {
	m := make(map[string]any)
	m["teacher_id"] = dto.TeacherID
	m["title"] = dto.Title
	m["description"] = dto.Description
	m["visibility"] = dto.Visibility

	if dto.StartTime != nil {
		m["start_time"] = *dto.StartTime
	}
	if dto.EndTime != nil {
		m["end_time"] = *dto.EndTime
	}

	query, args := r.qb.
		Insert("courses").
		SetMap(m).
		Suffix("RETURNING *").
		MustSql()

	var course Course
	err := r.storage.GetContext(ctx, &course, query, args...)
	if err != nil {
		return domain.Course{}, fmt.Errorf("failed to create course: %w", err)
	}

	return course.ToDomain(), nil
}

func (r *courseRepo) Delete(ctx context.Context, courseID string) (domain.Course, error) {
	query, args := r.qb.
		Delete("courses").
		Where(sq.Eq{"course_id": courseID}).
		Suffix("RETURNING *").
		MustSql()

	var course Course
	err := r.storage.GetContext(ctx, &course, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Course{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Course{}, fmt.Errorf("failed to delete course: %w", err)
	}

	return course.ToDomain(), nil
}

func (r *courseRepo) GetByID(ctx context.Context, courseID string) (domain.Course, error) {
	query, args := r.qb.
		Select("*").
		From("courses").
		Where(sq.Eq{"course_id": courseID}).
		MustSql()

	var course Course
	err := r.storage.GetContext(ctx, &course, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Course{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Course{}, fmt.Errorf("failed to get course: %v", err)
	}

	return course.ToDomain(), nil
}

func (r *courseRepo) ListByStudentID(ctx context.Context, studentID string) ([]domain.Course, error) {
	query, args := r.qb.
		Select("c.course_id", "c.teacher_id", "c.title", "c.description", "c.visibility", "c.start_time", "c.end_time", "c.created_at").
		From("enrollments e").
		Join("courses c ON e.course_id = c.course_id").
		Where(sq.Eq{"e.student_id": studentID}).
		MustSql()

	var courses []Course
	err := r.storage.SelectContext(ctx, &courses, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.Course{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %v", err)
	}

	res := make([]domain.Course, len(courses))
	for i, c := range courses {
		res[i] = c.ToDomain()
	}
	return res, nil
}

func (r *courseRepo) ListByTeacherID(ctx context.Context, teacherID string) ([]domain.Course, error) {
	query, args := r.qb.
		Select("course_id", "teacher_id", "title", "description", "visibility", "start_time", "end_time", "created_at").
		From("courses").
		Where(sq.Eq{"teacher_id": teacherID}).
		MustSql()

	var courses []Course
	err := r.storage.SelectContext(ctx, &courses, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.Course{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %v", err)
	}

	res := make([]domain.Course, len(courses))
	for i, c := range courses {
		res[i] = c.ToDomain()
	}
	return res, nil
}

func (r *courseRepo) Update(ctx context.Context, dto dto.UpdateCourseDTO) (domain.Course, error) {
	m := make(map[string]any)
	if dto.Title != nil {
		m["title"] = *dto.Title
	}
	if dto.Description != nil {
		m["description"] = *dto.Description
	}
	if dto.Visibility != nil {
		m["visibility"] = *dto.Visibility
	}
	if dto.StartTime != nil {
		m["start_time"] = *dto.StartTime
	}
	if dto.EndTime != nil {
		m["end_time"] = *dto.EndTime
	}
	query, args := r.qb.
		Update("courses").
		SetMap(m).
		Where(sq.Eq{"course_id": dto.ID}).
		Suffix("RETURNING *").
		MustSql()

	var course Course
	err := r.storage.GetContext(ctx, &course, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Course{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Course{}, fmt.Errorf("failed to update course: %v", err)
	}

	return course.ToDomain(), nil
}

// TODO fix return on conflict
func (r *courseRepo) EnrollUser(ctx context.Context, courseID, studentID string) (domain.Enrollment, error) {
	query, args := r.qb.
		Insert("enrollments").
		Columns("course_id", "student_id").
		Values(courseID, studentID).
		Suffix("ON CONFLICT DO NOTHING RETURNING *").
		MustSql()

	var enrollment Enrollment
	err := r.storage.GetContext(ctx, &enrollment, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		query, args = r.qb.
			Select("*").
			From("enrollments").
			Where(sq.Eq{"course_id": courseID, "student_id": studentID}).
			MustSql()
		err = r.storage.GetContext(ctx, &enrollment, query, args...)
	}
	if err != nil {
		return domain.Enrollment{}, fmt.Errorf("failed to enroll user: %v", err)
	}
	return enrollment.ToDomain(), nil
}

func (r *courseRepo) ExpelUser(ctx context.Context, courseID, studentID string) (domain.Enrollment, error) {
	query, args := r.qb.
		Delete("enrollments").
		Where(sq.Eq{"course_id": courseID, "student_id": studentID}).
		Suffix("RETURNING *").
		MustSql()

	var enrollment Enrollment
	err := r.storage.GetContext(ctx, &enrollment, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Enrollment{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Enrollment{}, fmt.Errorf("failed to expel user: %v", err)
	}

	return enrollment.ToDomain(), nil
}

func (r *courseRepo) IsTeacher(ctx context.Context, courseID, teacherID string) (bool, error) {
	query, args := r.qb.
		Select("TRUE").
		From("courses").
		Where(sq.Eq{"teacher_id": teacherID, "course_id": courseID}).
		MustSql()

	var isTeacher bool
	err := r.storage.GetContext(ctx, &isTeacher, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to verify teacher: %v", err)
	}

	return isTeacher, nil
}

func (r *courseRepo) IsMember(ctx context.Context, courseID, userID string) (bool, error) {
	query, args := r.qb.
		Select("TRUE").
		From("enrollments").
		Where(sq.Eq{"course_id": courseID, "student_id": userID}).
		MustSql()

	var isMember bool
	err := r.storage.GetContext(ctx, &isMember, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to verify member: %v", err)
	}

	return isMember, nil
}

func (r *courseRepo) ListCourseStudents(ctx context.Context, courseID string, index, limit int32) ([]domain.Student, int32, error) {
	query, args := r.qb.
		Select("u.user_id", "u.email", "u.first_name", "u.last_name").
		From("enrollments e").
		Join("users u ON u.user_id = e.student_id").
		Where(sq.Eq{"e.course_id": courseID}).
		Limit(uint64(limit)).
		Offset(uint64(index * limit)).
		MustSql()

	var members []Student
	err := r.storage.SelectContext(ctx, &members, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.Student{}, 0, nil
	}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list course students: %v", err)
	}

	res := make([]domain.Student, len(members))
	for i, m := range members {
		res[i] = m.ToDomain()
	}

	countQuery, countArgs := r.qb.
		Select("COUNT(*)").
		From("enrollments").
		Where(sq.Eq{"course_id": courseID}).
		MustSql()

	var total int32
	err = r.storage.GetContext(ctx, &total, countQuery, countArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count course students: %v", err)
	}

	return res, total, nil
}
