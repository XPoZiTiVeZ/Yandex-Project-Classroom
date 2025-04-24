package repo

import (
	"Classroom/Notifications/internal/domain"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type taskRepo struct {
	storage *sqlx.DB
	qb      sq.StatementBuilderType
}

func NewTaskRepo(storage *sqlx.DB) *taskRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &taskRepo{
		storage: storage,
		qb:      qb,
	}
}

func (r *taskRepo) GetByID(ctx context.Context, id string) (domain.Task, error) {
	query, args := r.qb.
		Select("task_id", "title").
		From("tasks").
		Where(sq.Eq{"task_id": id}).
		MustSql()

	var task Task
	err := r.storage.GetContext(ctx, &task, query, args...)
	if err != nil {
		return domain.Task{}, err
	}
	return task.ToDomain(), nil
}
