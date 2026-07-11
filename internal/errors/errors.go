package errors

import "errors"

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
