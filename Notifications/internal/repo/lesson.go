package repo

import (
	"Classroom/Notifications/internal/domain"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type lessonRepo struct {
	storage *sqlx.DB
	qb      sq.StatementBuilderType
}

func NewLessonRepo(storage *sqlx.DB) *lessonRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &lessonRepo{
		storage: storage,
		qb:      qb,
	}
}

func (r *lessonRepo) GetByID(ctx context.Context, id string) (domain.Lesson, error) {
	query, args := r.qb.
		Select("lesson_id", "title").
		From("lessons").
		Where(sq.Eq{"lesson_id": id}).
		MustSql()

	var lesson Lesson
	err := r.storage.GetContext(ctx, &lesson, query, args...)
	if err != nil {
		return domain.Lesson{}, err
	}
	return lesson.ToDomain(), nil
}
