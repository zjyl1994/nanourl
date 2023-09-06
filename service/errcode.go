package service

import "errors"

var (
	ErrNotFound  = errors.New("id not found")
	ErrDuplicate = errors.New("code duplicate")
)
