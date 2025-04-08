package repo

import (
	"Classroom/Lessons/internal/domain"
	"Classroom/Lessons/internal/dto"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type taskRepo struct {
	storage *sqlx.DB
	qb      sq.StatementBuilderType // Query Builder для удобного составления запросов
}

func NewTaskRepo(storage *sqlx.DB) *taskRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &taskRepo{
		storage: storage,
		qb:      qb,
	}
}

func (r *taskRepo) Create(ctx context.Context, payload dto.CreateTaskDTO) (domain.Task, error) {
	query, args := r.qb.
		Insert("tasks").
		Columns("course_id", "title", "content").
		Values(payload.CourseID, payload.Title, payload.Content).
		Suffix("RETURNING *").
		MustSql()

	var task Task
	if err := r.storage.GetContext(ctx, &task, query, args...); err != nil {
		return domain.Task{}, err
	}

	return task.ToEntity(), nil
}
