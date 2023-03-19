package shdb

import (
	"errors"
)

// Todo: Make it possible to add arguments to the errors.
// For instance NewErrNotFound(item string) error and then
// also a new IsErrNotFound(err).

var (
	ErrNotAnObject      = errors.New("not an object type")
	ErrInvalidType      = errors.New("invalid type")
	ErrNotFound         = errors.New("not found")
	ErrSessionInvalid   = errors.New("session invalid")
	ErrContextCancelled = errors.New("context cancelled")
	errJson             = errors.New("invalid json data")
)
