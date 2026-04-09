package errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInternalServer     = errors.New("internal server error")
)
