package domain

import "errors"

var (

	// General Errors
	ErrInternalServerError = errors.New("internal server error")
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUnexpected   = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrIncorrectUserID = errors.New("incorrect userID")
	ErrNoUpdate = errors.New("no update has been applied")
)