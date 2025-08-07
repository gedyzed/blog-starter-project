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
	ErrIncorrectUserID = errors.New("incorrect userID")
	ErrNoUpdate = errors.New("no update has been applied")
	ErrBadRequest = errors.New("invalid or bad request")
	ErrDuplicateKey = errors.New("duplicate key found")
	ErrInvalidUserID         = errors.New("invalid userID")

	// Token errors
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
