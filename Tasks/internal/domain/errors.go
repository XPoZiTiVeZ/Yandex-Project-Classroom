package domain

import "errors"

// Общие ошибки домена
var (
	ErrNotFound      = errors.New("entity not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyExists = errors.New("entity already exists")
	ErrDatabase      = errors.New("database error")
)
