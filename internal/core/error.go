package core

import (
	"errors"
)

var (
	// ErrNotFound represents a resource is not found.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists represents a resource already exists.
	ErrAlreadyExists = errors.New("already exists")
)
