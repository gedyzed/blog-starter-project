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
)

type UserUsecases struct {
	userRepo        domain.IUserRepository
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
		case domain.ErrUserNotFound:
			return nil, ErrInvalidCredential
		default:
			return nil, domain.ErrInternalServerError
		}
	}

	if err = u.passwordService.Verify(user.Password, data.Password); err != nil {
		return nil, ErrInvalidCredential
	}

	if data.Email != user.Email {
		return nil, ErrInvalidCredential
	}

	id := user.ID.Hex()
	token, err := u.tokenUsecase.GenerateTokens(ctx, id)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	return token, nil
}

func (u *UserUsecases) Authenticate(ctx context.Context, token string) (*domain.User, error) {
	userID, err := u.tokenUsecase.VerifyAccessToken(token)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	user, err := u.userRepo.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}

		return nil, domain.ErrInternalServerError
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
	if existing != nil {
		return domain.ErrUsernameAlreadyExists
	}

	// check if email exists
	existing, err = u.userRepo.GetByEmail(ctx, user.Email)
	if existing != nil {
		return domain.ErrEmailAlreadyExists
	}

	// hash the user password
	user.Password, err = u.passwordService.Hash(user.Password)
	if err != nil {
		return domain.ErrInternalServerError
	}

	// Add user to database
	return u.userRepo.Add(ctx, user)
}

func (u *UserUsecases) VerifyCode(ctx context.Context, vcode string) (string, error) {
	return u.tokenUsecase.VerifyCode(ctx, vcode)
}

func (u *UserUsecases) DeleteVCode(ctx context.Context, userID string) error {
	return u.tokenUsecase.DeleteVCode(ctx, userID)
}

func (u *UserUsecases) ForgotPassword(ctx context.Context, email string) error {

	// check if a user already exist
	_, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(domain.ErrInternalServerError, err) {
			return domain.ErrInternalServerError
		}

		return domain.ErrUserNotFound
	}


	return u.tokenUsecase.CreateSendVCode(ctx, email, Password_Reset)
}

func (u *UserUsecases) ResetPassword(ctx context.Context, email string, password string) error {

	password, err := u.passwordService.Hash(password)
	if err != nil {
		return domain.ErrInternalServerError
	}

	return u.userRepo.Update(ctx, "email", email, &domain.User{Password: password})
}

func (u *UserUsecases) PromoteDemote(ctx context.Context, userID string) error {

	existing, err := u.userRepo.Get(ctx, userID)
	if err != nil {
		return domain.ErrIncorrectUserID
	}

	user := &domain.User{}
	if existing.Role == "admin" {
		user.Role = "user"
	} else {
		user.Role = "admin"
	}

	return u.userRepo.Update(ctx, "_id", userID, user)
}

func (u *UserUsecases) ProfileUpdate(ctx context.Context, profileUpdate *domain.ProfileUpdateInput) error {

	existing, err := u.userRepo.Get(ctx, profileUpdate.UserID) 
	if err != nil {
		return err
	}

	existingProfile := &domain.ProfileUpdateInput{
		Firstname: existing.Firstname,
		Lastname: existing.Lastname,
		Profile: existing.Profile,
	} 

	if existingProfile == profileUpdate {
		return domain.ErrNoUpdate
	}

	user := &domain.User{
		Firstname: profileUpdate.Firstname,
		Lastname: profileUpdate.Lastname,
		Profile: profileUpdate.Profile,
	}

	return u.userRepo.Update(ctx, "_id", profileUpdate.UserID, user)
	
}
