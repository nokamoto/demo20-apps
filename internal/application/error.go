package application

import (
	"errors"
	"fmt"

	"github.com/nokamoto/demo20-apps/internal/mysql"
)

var (
	// ErrNotFound represents a resource is not found.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists represents a resource already exists.
	ErrAlreadyExists = errors.New("already exists")
	// ErrInternal represents an unexpected condition. It recommends to print the error message if caught.
	ErrInternal = errors.New("internal")
)

// ErrorCode is an internal error number for ErrorMap.
type ErrorCode int

// ErrorMap maps ErrorCode to a resource id.
type ErrorMap map[ErrorCode]string

const (
	// NotFound is ErrorCode corresponding to ErrNotFound.
	NotFound ErrorCode = iota
	// AlreadyExists is ErrorCode corresponing to ErrAlreadyExists.
	AlreadyExists
	// Internal is ErrorCode corresponding to ErrInternal.
	Internal
)

// Error converts the mysql error to an application error using the error map.
func Error(err error, m ErrorMap) error {
	res := func(code ErrorCode, unwrapped error) error {
		if id, ok := m[code]; ok {
			return fmt.Errorf("%s: %w", id, unwrapped)
		}
		return unwrapped
	}

	if errors.Is(err, mysql.ErrNotFound) {
		return res(NotFound, ErrNotFound)
	}

	if errors.Is(err, mysql.ErrAlreadyExists) {
		return res(AlreadyExists, ErrAlreadyExists)
	}

	return err
}
