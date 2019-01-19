package manager

import "errors"

var (
	errNotExist     = errors.New("resource not found")
	errExist        = errors.New("resource exist")
	errPermission   = errors.New("permission denied")
	errInvalidInput = errors.New("invalid input")
)
