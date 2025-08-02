package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gedyzed/blog-starter-project/Repository"
)

var (
	// Access token errors
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrExpiredAccessToken = errors.New("access token has expired")

	// Refresh token errors
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrExpiredRefreshToken = errors.New("refresh token has expired")

	// input errors
	ErrInvalidCredential = errors.New("invalid email or passwrod")

	// General errors
	ErrUnexpected   = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized")
)

type UserUsecases struct {
	tokenRepo       domain.ITokenRepository
	userRepo        domain.IUserRepository
	tokenService    domain.ITokenService
	passwordService domain.IPasswordService
}

func NewUserUsecase(repo domain.IUserRepository, ts domain.ITokenService, ps domain.IPasswordService) *UserUsecases {
	return &UserUsecases{
		userRepo:        repo,
		tokenService:    ts,
		passwordService: ps,
	}
}

func (u *UserUsecases) Login(ctx context.Context, user domain.User) (*domain.Token, error) {
	data, err := u.userRepo.Get(ctx, user.Username)

	if err != nil {
		switch err {
		case repository.ErrUserNotFound:
			return nil, ErrInvalidCredential
		default:
			return nil, ErrUnexpected
		}
	}

	if err = u.passwordService.Verify(user.Password, data.Password); err != nil {
		return nil, errors.New("invalid email or passwrod")
	}

	if data.Email != user.Email {
		return nil, ErrInvalidCredential
	}

	token, err := u.tokenService.GenerateTokens(ctx, data.ID)
	if err != nil {
		return nil, ErrUnexpected
	}

	return token, nil
}

func (u *UserUsecases) Authenticate(ctx context.Context, token string) (*domain.User, error) {
	userID, err := u.tokenService.VerifyAccessToken(token)
	if err != nil {
		return nil, ErrUnauthorized
	}

	user, err := u.userRepo.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUnauthorized
		}

		return nil, ErrUnexpected
	}

	return user, nil
}

func (u *UserUsecases) RefreshToken(ctx context.Context, id string, refreshToken string) (*domain.Token, error) {
	token, err := u.tokenRepo.FindByUserID(ctx, id)
	if err != nil {
		return nil, nil
	}

	if token.RefreshExpiry.Unix() > time.Now().Unix() {
		err = u.tokenRepo.DeleteByUserID(ctx, id)
		if err != nil {
			return nil, nil
		}
		return nil, ErrExpiredRefreshToken
	}

	return u.tokenService.RefreshTokens(ctx, refreshToken)
}

func (u *UserUsecases) Register(ctx context.Context, user *domain.User) error {

	// check if username exists
	existing, err := u.userRepo.GetByUsername(ctx, user.Username)
	if err != nil {
		err = checkAndReturnError(existing, err, "username already exists")
		if err != nil {
			return err
		}
	}

	// check if email exists
	existing, err = u.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		err = checkAndReturnError(existing, err, "email already exists")
		if err != nil {
			return err
		}
	}

	// hash the user password
	user.Password, err = u.passwordService.Hash(user.Password)
	if err != nil {
		return errors.New("internal server error")
	}

	return u.userRepo.Add(ctx, user)
}

func checkAndReturnError(existing *domain.User, err error, errorMessage string) error {

	if err != nil {
		if err.Error() == "error while decoding data" {
			return errors.New("internal server error")
		} else if err.Error() == "internal server error" {
			return errors.New("internal server error")
		}
	}

	if existing != nil {
		return errors.New(errorMessage)
	}

	return nil

}
