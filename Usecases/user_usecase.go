package usecases

import (
	"errors"

	"github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gedyzed/blog-starter-project/Infrastructure"
)

type UserUsecases struct {
	repo            domain.IUserRepository
	tokenService    domain.ITokenService
	passwordService domain.IPasswordService
}

func (u *UserUsecases) Login(user domain.User) (*domain.Token, error) {
	data, err := u.repo.Get(user.Username)
	if err != nil {
		return nil, errors.New("User does not exist")
	}

	if u.passwordService.Verify(user.Password, data.Password) {
		return nil, errors.New("invalid email or passwrod")
	}

	if data.Email != user.Email {
		return nil, errors.New("invalid email or passwrod")
	}

	token, err := u.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (u *UserUsecases) Authenticate(token string) error {
	return u.tokenService.ValidateToken(token)
}

func (u *UserUsecases) RefreshToken(refreshToken string) (string, error) {
	return u.tokenService.RefreshToken(refreshToken)
}
