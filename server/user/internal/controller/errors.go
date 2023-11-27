package controller

import "errors"

var ErrTokenInvalid = errors.New("token invalid")
var ErrTokenExpired = errors.New("token expired")
