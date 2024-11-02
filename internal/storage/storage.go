package storage

import "errors"

var (
	ErrTaskAlreadyExists = errors.New("task already exists")
	ErrTaskNotFound      = errors.New("task not found")
	ErrAppNotFound       = errors.New("app not found")
)
