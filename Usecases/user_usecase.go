package usecases

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gedyzed/blog-starter-project/Domain"
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
	ErrUserNotFound      = errors.New("user not found")
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
	data, err := u.userRepo.GetByUsername(ctx, user.Username)

	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			return nil, err
		default:
			log.Println(err.Error())
			return nil, domain.ErrInternalServer
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
		log.Println(err.Error())
		return nil, domain.ErrInternalServer
	}

	return token, nil
}

func (u *UserUsecases) Authenticate(ctx context.Context, token string) (*domain.User, error) {
	userID, err := u.tokenUsecase.VerifyAccessToken(token)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}

		return nil, domain.ErrInternalServer
	}

	return user, nil
}

func (u *UserUsecases) RefreshToken(ctx context.Context, id string, refreshToken string) (*domain.Token, error) {

	token, err := u.tokenUsecase.FindByUserID(ctx, id)
	if err != nil {
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

func (u *UserUsecases) Register(ctx context.Context, user *domain.User) (string, error) {

	// Ensure Provider is set
	if user.Provider == "" {
		user.Provider = "local"
	}

	// Ensure User Role 
	if user.Role != "user"{
		user.Role = "user"
	}

	// Check email uniqueness
	existing, err := u.userRepo.GetByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return "", domain.ErrInternalServer
	}
	if existing != nil {
		return "", domain.ErrEmailAlreadyExists
	}


	// Check username uniqueness (if provided)
	if user.Username != "" {
		existing, err = u.userRepo.GetByUsername(ctx, user.Username)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInternalServer
		}
		if existing != nil {
			return "", domain.ErrUsernameAlreadyExists
		}
	}

	// Handle password
	if user.Provider == "local" {
		user.Password, err = u.passwordService.Hash(user.Password)
		if err != nil {
			return "", domain.ErrInternalServer
		}
	} else {
		user.Password = ""
	}

	// Set timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save user to DB
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
		if errors.Is(err, domain.ErrInternalServer) {
			return domain.ErrInternalServer
		}

		return domain.ErrUserNotFound
	}

	return u.tokenUsecase.CreateSendVCode(ctx, email, Password_Reset)
}

func (u *UserUsecases) ResetPassword(ctx context.Context, email string, password string) error {

	password, err := u.passwordService.Hash(password)
	if err != nil {
		return domain.ErrInternalServer
	}

	return u.userRepo.Update(ctx, "email", email, &domain.User{Password: password})
}

func (u *UserUsecases) PromoteDemote(ctx context.Context, userID string) error {

	existing, err := u.userRepo.Get(ctx, userID)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			return err
		default:
			return domain.ErrInternalServer
		}
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

	if (profileUpdate.Firstname == "" || profileUpdate.Firstname == existing.Firstname) &&
		(profileUpdate.Lastname == "" || profileUpdate.Lastname == existing.Lastname) &&
		(profileUpdate.Bio == "" || profileUpdate.Bio == existing.Profile.Bio) &&
		(profileUpdate.ProfilePic == "" || profileUpdate.ProfilePic == existing.Profile.ProfilePic) &&
		(profileUpdate.Location == "" || profileUpdate.Location == existing.Profile.ContactInfo.Location) &&
		(profileUpdate.PhoneNumber == "" || profileUpdate.PhoneNumber == existing.Profile.ContactInfo.PhoneNumber) {
		return domain.ErrNoUpdate
	}

	user := &domain.User{
		Firstname: profileUpdate.Firstname,
		Lastname: profileUpdate.Lastname,
		Profile: domain.Profile{
			Bio: profileUpdate.Bio,
			ContactInfo: domain.ContactInformation{
				Location: profileUpdate.Location,
				PhoneNumber: profileUpdate.PhoneNumber,
			},
			ProfilePic: profileUpdate.ProfilePic,
		},
	}

	return u.userRepo.Update(ctx, "_id", profileUpdate.UserID, user)
}

func (u *UserUsecases) SaveToken (ctx context.Context, tokens *domain.Token) error {
	return u.tokenUsecase.SaveToken(ctx, tokens)
}



