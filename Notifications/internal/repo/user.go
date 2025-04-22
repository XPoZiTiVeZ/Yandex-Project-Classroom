package repo

import (
	"Classroom/Notifications/internal/domain"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	storage *sqlx.DB
	qb      sq.StatementBuilderType
}

func NewUserRepo(storage *sqlx.DB) *userRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &userRepo{
		storage: storage,
		qb:      qb,
	}
}

func (r *userRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	query, args := r.qb.
		Select("user_id", "email", "first_name", "last_name").
		From("users").
		Where(sq.Eq{"user_id": id}).
		MustSql()

	var user User
	err := r.storage.GetContext(ctx, &user, query, args...)
	if err != nil {
		return domain.User{}, err
	}
	return user.ToDomain(), nil
}

func (r *userRepo) ListByCourseID(ctx context.Context, courseID string) ([]domain.User, error) {
	query, args := r.qb.
		Select("u.user_id", "u.email", "u.first_name", "u.last_name").
		From("enrollments e").
		Join("users u ON u.user_id = e.student_id").
		Where(sq.Eq{"e.course_id": courseID}).
		MustSql()

	var users []User
	err := r.storage.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, err
	}

	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = user.ToDomain()
	}

	return domainUsers, nil
}
