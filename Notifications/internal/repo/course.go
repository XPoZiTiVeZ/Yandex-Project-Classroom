package repo

import (
	"Classroom/Notifications/internal/domain"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type courseRepo struct {
	storage *sqlx.DB
	qb      sq.StatementBuilderType
}

func NewCourseRepo(storage *sqlx.DB) *courseRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &courseRepo{
		storage: storage,
		qb:      qb,
	}
}

func (r *courseRepo) GetByID(ctx context.Context, id string) (domain.Course, error) {
	query, args := r.qb.
		Select("course_id", "title").
		From("courses").
		Where(sq.Eq{"course_id": id}).
		MustSql()

	var course Course
	err := r.storage.GetContext(ctx, &course, query, args...)
	if err != nil {
		return domain.Course{}, err
	}
	return course.ToDomain(), nil
}
