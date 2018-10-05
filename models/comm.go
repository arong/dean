package models

import "errors"

var (
	errNotExist = errors.New("resource not found")
	errExist = errors.New("resource exist")
)
