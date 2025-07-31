package usecases

import (
	"context"
	"errors"

	"github.com/gedyzed/blog-starter-project/Domain"
)

type UserUsecases struct {
	repo            domain.IUserRepository
	tokenService    domain.ITokenService
	passwordService domain.IPasswordService
}

func NewUserUsecase(repo domain.IUserRepository, ts domain.ITokenService, ps domain.IPasswordService) *UserUsecases{
	return &UserUsecases{
		repo: repo, 
		tokenService: ts, 
		passwordService: ps,
	}
}

func (u *UserUsecases) Login(ctx context.Context, user domain.User) (*domain.Token, error) {

	data, err := u.repo.Get(ctx, user.Username)
	if err != nil {
		return nil, errors.New("user does not exist")
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


func (u *UserUsecases) Register(ctx context.Context, user *domain.User) error {

	// check if username exists
	existing, err := u.repo.GetByUsername(ctx, user.Username)
	if err != nil {
		err = checkAndReturnError(existing, err, "username already exists")
		if err != nil{
			return err
		}	 
	}

	// check if email exists 
	existing, err = u.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		err = checkAndReturnError(existing, err, "email already exists")
		if err != nil{
			return err
		}	 
	}

	// hash the user password
	user.Password, err = u.passwordService.Hash(user.Password)
	if err != nil {
		return errors.New("internal server error")
	}

	return u.repo.Add(ctx, user)
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
