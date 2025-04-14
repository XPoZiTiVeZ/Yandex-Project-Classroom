package repo

import (
	"Classroom/Tasks/internal/domain"
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type statusesRepo struct {
	storage *sqlx.DB
	qb      sq.StatementBuilderType // Query Builder для удобного составления запросов
}

func NewStatusesRepo(storage *sqlx.DB) *statusesRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &statusesRepo{
		storage: storage,
		qb:      qb,
	}
}

func (r *statusesRepo) Get(ctx context.Context, taskID, userID string) (domain.TaskStatus, error) {
	query, args := r.qb.
		Select("*").
		From("task_submissions").
		Where(sq.Eq{"task_id": taskID, "student_id": userID}).
		MustSql()

	var status TaskStatus
	err := r.storage.GetContext(ctx, &status, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.TaskStatus{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.TaskStatus{}, err
	}
	return status.ToEntity(), nil
}

func (r *statusesRepo) Update(ctx context.Context, status domain.TaskStatus) error {
	query, args := r.qb.
		Update("task_submissions").
		Set("completed", status.Completed).
		Where(sq.Eq{"task_id": status.TaskID, "student_id": status.UserID}).
		MustSql()

	_, err := r.storage.ExecContext(ctx, query, args...)
	return err
}

func (r *statusesRepo) Create(ctx context.Context, status domain.TaskStatus) error {
	query, args := r.qb.
		Insert("task_submissions").
		Columns("task_id", "student_id", "completed").
		Values(status.TaskID, status.UserID, status.Completed).
		MustSql()

	_, err := r.storage.ExecContext(ctx, query, args...)
	return err
}

func (r *statusesRepo) ListByTaskID(ctx context.Context, taskID string) ([]domain.TaskStatus, error) {
	query, args := r.qb.
		Select(
			"t.task_id",
			"e.student_id",
			"COALESCE(ts.completed, FALSE) AS completed",
		).
		From("tasks t").
		Join("enrollments e ON e.course_id = t.course_id").
		LeftJoin("task_submissions ts ON ts.task_id = t.task_id AND ts.student_id = e.student_id").
		Where(sq.Eq{"t.task_id": taskID}).
		MustSql()

	var statuses []TaskStatus
	err := r.storage.SelectContext(ctx, &statuses, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.TaskStatus{}, nil
	}
	if err != nil {
		return nil, err
	}

	result := make([]domain.TaskStatus, len(statuses))
	for i, status := range statuses {
		result[i] = status.ToEntity()
	}

	return result, nil
}
