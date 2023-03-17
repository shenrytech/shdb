package shdb

import "errors"

var (
	ErrNotAnObject = errors.New("not an object type")
	ErrInvalidType = errors.New("invalid type")
	ErrNotFound    = errors.New("not found")
)
