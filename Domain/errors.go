package domain

import (
	"errors"
)

var (

	// General Errors
	ErrInternalServer        = errors.New("internal server error")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExist      = errors.New("user already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUnexpected            = errors.New("internal server error")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidUserID         = errors.New("invalid userID")
	ErrNoUpdate              = errors.New("no update has been applied")

	// Token errors
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
