package controller

import "errors"

var (
	ErrTokenInvalid = errors.New("token invalid")
	ErrTokenExpired = errors.New("token expired")
	ErrNotFound = errors.New("not found")
)
