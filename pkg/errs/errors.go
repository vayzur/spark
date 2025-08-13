package errs

import "errors"

var (
	ErrTagExists = errors.New("tag exists")
	ErrNotFound  = errors.New("not found")
)
