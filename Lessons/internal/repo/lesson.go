package repo

import (
	"Classroom/Lessons/internal/domain"
	"Classroom/Lessons/internal/dto"
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type lessonRepo struct {
	storage *sqlx.DB
	qb      sq.StatementBuilderType // Query Builder для удобного составления запросов
}

func NewLessonRepo(storage *sqlx.DB) *lessonRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &lessonRepo{
		storage: storage,
		qb:      qb,
	}
}

func (r *lessonRepo) Create(ctx context.Context, dto dto.CreateLessonDTO) (domain.Lesson, error) {
	query, args := r.qb.
		Insert("lessons").
		Columns("course_id", "title", "content").
		Values(dto.CourseID, dto.Title, dto.Content).
		Suffix("RETURNING *").
		MustSql()

	var lesson Lesson
	if err := r.storage.GetContext(ctx, &lesson, query, args...); err != nil {
		return domain.Lesson{}, err
	}
	return lesson.ToEntity(), nil
}

func (r *lessonRepo) GetByID(ctx context.Context, id string) (domain.Lesson, error) {
	query, args := r.qb.
		Select("*").
		From("lessons").
		Where(sq.Eq{"lesson_id": id}).
		MustSql()

	var lesson Lesson
	err := r.storage.GetContext(ctx, &lesson, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Lesson{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Lesson{}, err
	}
	return lesson.ToEntity(), nil
}

func (r *lessonRepo) ListByCourseID(ctx context.Context, courseID string) ([]domain.Lesson, error) {
	query, args := r.qb.
		Select("*").
		From("lessons").
		Where(sq.Eq{"course_id": courseID}).
		MustSql()
	var lessons []Lesson
	err := r.storage.SelectContext(ctx, &lessons, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.Lesson{}, nil
	}
	if err != nil {
		return nil, err
	}
	var result []domain.Lesson
	for _, l := range lessons {
		result = append(result, l.ToEntity())
	}
	return result, nil
}

func (r *lessonRepo) Update(ctx context.Context, dto dto.UpdateLessonDTO) (domain.Lesson, error) {
	m := make(map[string]any)
	if dto.Title != nil {
		m["title"] = *dto.Title
	}
	if dto.Content != nil {
		m["content"] = *dto.Content
	}
	query, args := r.qb.
		Update("lessons").
		SetMap(m).
		Where(sq.Eq{"lesson_id": dto.LessonID}).
		Suffix("RETURNING *").
		MustSql()

	var lesson Lesson
	err := r.storage.GetContext(ctx, &lesson, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Lesson{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Lesson{}, err
	}
	return lesson.ToEntity(), nil
}

func (r *lessonRepo) Delete(ctx context.Context, id string) error {
	query, args := r.qb.
		Delete("lessons").
		Where(sq.Eq{"lesson_id": id}).
		MustSql()
	res, err := r.storage.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return domain.ErrNotFound
	}
	return nil
}
