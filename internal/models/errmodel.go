package models

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrUserNotFound = errors.New("user not found")
)
