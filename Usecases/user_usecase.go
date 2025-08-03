package usecases

import (
	"context"
	"errors"
	"time"
	

	domain "github.com/gedyzed/blog-starter-project/Domain"
	repository "github.com/gedyzed/blog-starter-project/Repository"
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
	userRepo  			domain.IUserRepository
	tokenUsecase    ITokenUsecase
	passwordService domain.IPasswordService
}

func NewUserUsecase(userRepo domain.IUserRepository, tu ITokenUsecase, ps domain.IPasswordService) *UserUsecases {
	return &UserUsecases{
		userRepo:        userRepo,
		tokenUsecase:    tu,
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
		return nil, ErrInvalidCredential
	}

	if data.Email != user.Email {
		return nil, ErrInvalidCredential
	}

	token, err := u.tokenUsecase.GenerateTokens(ctx, data.ID)
	if err != nil {
		return nil, ErrUnexpected
	}

	return token, nil
}

func (u *UserUsecases) Authenticate(ctx context.Context, token string) (*domain.User, error) {
	userID, err := u.tokenUsecase.VerifyAccessToken(token)
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

	token, err := u.tokenUsecase.FindByUserID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTokenNotFound) {
			return nil, ErrInvalidCredential
		}
		return nil, err
	}

	if token.RefreshExpiry.Unix() > time.Now().Unix() {
		err = u.tokenUsecase.DeleteByUserID(ctx, id)
		if err != nil {
			return nil, err
		}

		return nil, ErrExpiredRefreshToken
	}

	return u.tokenUsecase.RefreshTokens(ctx, refreshToken)
}


func (u *UserUsecases) Register(ctx context.Context, user *domain.User) error {

	// check if username exists
	existing, err := u.userRepo.GetByUsername(ctx, user.Username)
	if err != nil && err.Error() == "error while decoding data"{
		return errors.New("internal server error")
	}

	if existing != nil {
		return errors.New("username already exists")
	}

	// check if email exists
	existing, err = u.userRepo.GetByEmail(ctx, user.Email)
		if err != nil && err.Error() == "error while decoding data"{
		return errors.New("internal server error")
	}

	if existing != nil {
		return errors.New("email already exists")
	}

	// hash the user password
	user.Password, err = u.passwordService.Hash(user.Password)
	if err != nil {
		return errors.New("internal server error")
	}

	// Add user to database 
	return u.userRepo.Add(ctx, user)
}

func (u *UserUsecases) VerifyCode(ctx context.Context, userID string, vcode string) error {
	return u.tokenUsecase.VerifyCode(ctx, userID, vcode)
}

func (u *UserUsecases) DeleteVCode(ctx context.Context, userID string) error {
	return u.tokenUsecase.DeleteVCode(ctx, userID)
}


 