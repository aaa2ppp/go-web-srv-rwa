package dto

import "errors"

var (
	ErrUnauthorized  = errors.New("unauthorized")
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)
