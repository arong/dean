package models

import "errors"

var (
	errNotExist     = errors.New("resource not found")
	errExist        = errors.New("resource exist")
	errPermission   = errors.New("permission denied")
	errInvalidInput = errors.New("invalid input")
)

const (
	// TypeStudent => student
	TypeStudent = 1
	// TypeTeacher => teacher
	TypeTeacher = 2
)
