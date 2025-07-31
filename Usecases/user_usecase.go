package usecases

import (
	"errors"

	"github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gedyzed/blog-starter-project/Infrastructure"
)

type UserUsecases struct {
	repo         domain.IUserRepository
	tokenService infrastructure.ITokenService
}

func (u *UserUsecases) Login(user domain.User) (string, error) {
	data, err := u.repo.Get(user.Username)
	if err != nil {
		return "", errors.New("User does not exist")
	}

	if data.Password != user.Password || data.Email != user.Email {
		return "", errors.New("invalid email or passwrod")
	}

	token, err := u.tokenService.GenerateToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUsecases) Authenticate(token string) error {
	return u.tokenService.ValidateToken(token)
}

func (u *UserUsecases) RefreshToken(refreshToken string) (string, error) {
	return u.tokenService.RefreshToken(refreshToken)
}
