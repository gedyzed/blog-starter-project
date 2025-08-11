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
	ErrIncorrectUserID       = errors.New("incorrect userID")
	ErrIncorrectUserID       = errors.New("incorrect userID")
	ErrNoUpdate              = errors.New("no update has been applied")
	ErrBadRequest            = errors.New("invalid or bad request")
	ErrDuplicateKey          = errors.New("duplicate key found")
	ErrInvalidUserID         = errors.New("invalid userID")
	ErrBadRequest            = errors.New("invalid or bad request")
	ErrDuplicateKey          = errors.New("duplicate key found")
	ErrInvalidUserID         = errors.New("invalid userID")

	// Token errors
	ErrInvalidToken        		= errors.New("invalid access token")
	ErrInvalidRefreshToken 		= errors.New("invalid refresh token")
	ErrMissingOrInvalidHeader 	= errors.New("missing or invalid authorization header")
	ErrTokenDoesNotMatch		= errors.New("token does not match the stored token")
	ErrTokenNotFound            = errors.New("token not found")

	// OAuth errors 
	ErrFailedToDecodeUserInfo = errors.New("failed to decode user information")
	ErrInvalidGoogleID	      = errors.New("invalid Google OAuth2 token")
	ErrFailedToFetch  	  	  = errors.New("failed to refresh access token")
	ErrFailedToExchange    	  = errors.New("failed to exchange authorization code")
	ErrFailedToFetchUserInfo  = errors.New("failed to fetch user information from Google")

	//Email Errors
	ErrFailedToSendEmail = errors.New("failed to send email")

	
)
