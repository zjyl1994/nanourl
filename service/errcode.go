package service

import "errors"

var (
	ErrCodeNotFound  = errors.New("code not found")
	ErrCodeDuplicate = errors.New("code duplicate")
	ErrCodeExhausted  = errors.New("code exhausted")
)
