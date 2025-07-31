package usecases

import (
	"errors"
	"time"

	"github.com/gedyzed/blog-starter-project/Domain"
)

type UserUsecases struct {
	tokenRepo       domain.ITokenRepository
	userRepo        domain.IUserRepository
	tokenService    domain.ITokenService
	passwordService domain.IPasswordService
}

func (u *UserUsecases) Login(user domain.User) (*domain.Token, error) {
	data, err := u.userRepo.Get(user.Username)
	if err != nil {
		return nil, errors.New("User does not exist")
	}

	if err = u.passwordService.Verify(user.Password, data.Password); err != nil {
		return nil, errors.New("invalid email or passwrod")
	}

	if data.Email != user.Email {
		return nil, errors.New("invalid email or passwrod")
	}

	token, err := u.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (u *UserUsecases) Authenticate(token string) error {
	return u.tokenService.ValidateToken(token)
}

func (u *UserUsecases) RefreshToken(id string, refreshToken string) (string, error) {
	token, err := u.tokenRepo.GetTokenByUserID(id)
	if err != nil {
		return "", nil
	}

	if token.ExpiresAt.Unix() > time.Now().Unix() {
		err = u.tokenRepo.Delete(token.ID)
		if err != nil {
			return "", nil
		}
		return "", errors.New("refresh token expired")
	}

	return u.tokenService.RefreshToken(refreshToken)
}
