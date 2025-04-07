package repo

import (
	"Classroom/Auth/internal/dto"
	"Classroom/Auth/internal/entities"
	"Classroom/Auth/internal/service"
	"context"
	"errors"

	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (r *userRepo) Create(ctx context.Context, dto dto.CreateUserDTO) (entities.User, error) {
	query, args := r.qb.
		Insert("users").
		Columns("email", "password_hash", "is_superuser", "first_name", "last_name").
		Values(dto.Email, dto.PasswordHash, dto.IsSuperUser, dto.FirstName, dto.LastName).
		Suffix("RETURNING *").
		MustSql()

	var user User
	if err := r.storage.GetContext(ctx, &user, query, args...); err != nil {
		if isUniqueViolation(err) {
			return entities.User{}, service.ErrUserAlreadyExists
		}
		return entities.User{}, err
	}
	return user.ToEntity(), nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (entities.User, error) {
	query, args := r.qb.
		Select("user_id", "email", "password_hash", "is_superuser", "first_name", "last_name").
		From("users").
		Where(sq.Eq{"email": email}).
		MustSql()

	var user User
	err := r.storage.GetContext(ctx, &user, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return entities.User{}, service.ErrUserNotFound
	}
	if err != nil {
		return entities.User{}, err
	}
	return user.ToEntity(), nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (entities.User, error) {
	query, args := r.qb.
		Select("user_id", "email", "password_hash", "is_superuser", "first_name", "last_name").
		From("users").
		Where(sq.Eq{"user_id": id}).
		MustSql()

	var user User
	err := r.storage.GetContext(ctx, &user, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return entities.User{}, service.ErrUserNotFound
	}
	if err != nil {
		return entities.User{}, err
	}
	return user.ToEntity(), nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code.Name() == "unique_violation"
	}
	return false
}
