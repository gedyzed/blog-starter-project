package usecases

import (
	"context"
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

<<<<<<< HEAD
func (u *UserUsecases) Login(user domain.User) (*domain.Token, error) {
	data, err := u.userRepo.Get(user.Username)
=======
func NewUserUsecase(repo domain.IUserRepository, ts domain.ITokenService, ps domain.IPasswordService) *UserUsecases{
	return &UserUsecases{
		repo: repo, 
		tokenService: ts, 
		passwordService: ps,
	}
}

func (u *UserUsecases) Login(ctx context.Context, user domain.User) (*domain.Token, error) {

	data, err := u.repo.Get(ctx, user.Username)
>>>>>>> 5326082c22b972240493a44a880a1835f7d591f8
	if err != nil {
		return nil, errors.New("user does not exist")
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
