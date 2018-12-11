package models

import "errors"

var (
	errNotExist = errors.New("resource not found")
	errExist    = errors.New("resource exist")
)

const (
	// TypeStudent => student
	TypeStudent = 1
	// TypeTeacher => teacher
	TypeTeacher = 2
)
