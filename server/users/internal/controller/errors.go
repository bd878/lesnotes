package controller

import "errors"

var (
	ErrTokenExpired  = errors.New("token expired")
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
	ErrUserExists    = errors.New("user exists")
)
