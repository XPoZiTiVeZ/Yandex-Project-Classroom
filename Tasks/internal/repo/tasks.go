package repo

import (
	"Classroom/Tasks/internal/domain"
	"Classroom/Tasks/internal/dto"
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

func (r *taskRepo) ListByCourseID(ctx context.Context, courseID string) ([]domain.Task, error) {
	query, args := r.qb.
		Select("*").
		From("tasks").
		Where(sq.Eq{"course_id": courseID}).
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

func (r *taskRepo) Update(ctx context.Context, task domain.Task) error {
	query, args := r.qb.
		Update("tasks").
		Set("title", task.Title).
		Set("content", task.Content).
		Where(sq.Eq{"task_id": task.ID}).
		MustSql()

	_, err := r.storage.ExecContext(ctx, query, args...)
	return err
}

// Метод проверяет наличие записи в БД и удаляет ее
func (r *taskRepo) Delete(ctx context.Context, id string) error {
	query, args := r.qb.
		Delete("tasks").
		Where(sq.Eq{"task_id": id}).
		MustSql()

	res, err := r.storage.ExecContext(ctx, query, args...)
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return domain.ErrNotFound
	}
	return err
}

func (r *taskRepo) GetTaskStatus(ctx context.Context, taskID, userID string) (domain.TaskStatus, error) {
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

func (r *taskRepo) UpdateTaskStatus(ctx context.Context, status domain.TaskStatus) error {
	query, args := r.qb.
		Update("task_submissions").
		Set("completed", status.IsCompleted).
		Where(sq.Eq{"task_id": status.TaskID, "student_id": status.UserID}).
		MustSql()

	_, err := r.storage.ExecContext(ctx, query, args...)
	return err
}

func (r *taskRepo) CreateTaskStatus(ctx context.Context, status domain.TaskStatus) error {
	query, args := r.qb.
		Insert("task_submissions").
		Columns("task_id", "student_id", "completed").
		Values(status.TaskID, status.UserID, status.IsCompleted).
		MustSql()

	_, err := r.storage.ExecContext(ctx, query, args...)
	return err
}
