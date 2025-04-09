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

func (r *taskRepo) GetTaskByID(ctx context.Context, id string) (domain.Task, error) {
	query, args := r.qb.
		Select("*").
		From("tasks").
		Where(sq.Eq{"task_id": id}).
		MustSql()

	var task Task
	err := r.storage.GetContext(ctx, &task, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Task{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Task{}, err
	}

	return task.ToEntity(), nil
}

func (r *taskRepo) ListByCourseID(ctx context.Context, course_id string) ([]domain.Task, error) {
	query, args := r.qb.
		Select("*").
		From("tasks").
		Where(sq.Eq{"course_id": course_id}).
		MustSql()

	var tasks []Task
	err := r.storage.SelectContext(ctx, &tasks, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.Task{}, nil
	}
	if err != nil {
		return nil, err
	}

	result := make([]domain.Task, len(tasks))
	for i, task := range tasks {
		result[i] = task.ToEntity()
	}

	return result, nil
}
