package repo

import (
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
