package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourseRepo struct {
	db *pgxpool.Pool
}

func NewCourseRepo(db *pgxpool.Pool) *CourseRepo {
	return &CourseRepo{db: db}
}
