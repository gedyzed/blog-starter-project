package usecases

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	domain "github.com/gedyzed/blog-starter-project/Domain"
)

var (
	Email_Verification = "email_verification"
	Password_Reset     = "password_reset"

	ResetPasswordEmailSubject  = "Sending Password Reset Link"
	ResetPasswordEmailBodyText = "Here is the link to reset your password click the link "
	ResetPasswordRoute         = "/users/reset-password?token="

	EmailVerificationSubject = "Sending Email Verification Code"
	EmailVerificationBody    = "Here is you verification code: "
)

var (
	ErrIncorrectToken = errors.New("incorrect token")
	ErrExpiredToken   = errors.New("expired token")
)

type ITokenUsecase interface {
	CreateSendVCode(ctx context.Context, userID string, tokenType string) error
	GenerateSecureToken(string) (string, error)
	VerifyCode(ctx context.Context, token *domain.VToken)(string, error)
	DeleteVCode(ctx context.Context, userID string) error
	FindByUserID(ctx context.Context, userID string) (*domain.Token, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*domain.Token, error)
	GenerateTokens(ctx context.Context, userID string) (*domain.Token, error)
	VerifyAccessToken(string) (string, error)
	DeleteByUserID(ctx context.Context, email string) error	
	SaveToken(ctx context.Context, token *domain.Token) error
	GetByAccessToken(ctx context.Context, accessToken string)(string, error)
	

}

type tokenUsecase struct {
	tokenRepo      domain.ITokenRepo
	vtokenRepo     domain.IVTokenRepo
	vtokenServices domain.IVTokenService
	tokenService   domain.ITokenService
}

func NewTokenUsecase(tokenRepo domain.ITokenRepo, vtokenRepo domain.IVTokenRepo, svs domain.IVTokenService, js domain.ITokenService) ITokenUsecase {
	return &tokenUsecase{
		tokenRepo:      tokenRepo,
		vtokenRepo:     vtokenRepo,
		vtokenServices: svs,
		tokenService:   js,
	}
}

func (t *tokenUsecase) CreateSendVCode(ctx context.Context, email string, tokenType string) error {

	// generate random verfication code
	token, err := t.GenerateSecureToken(tokenType)
	if err != nil {
		return err
	}

	// ten minutes of expiration time
	expiration_time := time.Now().Add(10 * time.Minute)

	vtoken := domain.VToken{
		Email:     email,
		TokenType: tokenType,
		Token:     token,
		ExpiresAt: expiration_time,
	}

	// save the created verification code to db
	err = t.vtokenRepo.CreateVCode(ctx, &vtoken)
	if err != nil {
		return domain.ErrInternalServer
	}

	if vtoken.TokenType == Email_Verification {
		return t.vtokenServices.SendEmail(
			[]string{email},
			EmailVerificationSubject,
			EmailVerificationBody + token,
		)
	}

	return t.vtokenServices.SendEmail(
		[]string{email},
		ResetPasswordEmailSubject,
		ResetPasswordRoute+token,
	)
}

func (t *tokenUsecase) GenerateSecureToken(tokenType string) (string, error) {

	if tokenType == Password_Reset {
		return rand.Text(), nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func (t *tokenUsecase) VerifyCode(ctx context.Context, token *domain.VToken) (string, error) {
	var (
		existingToken *domain.VToken
		err           error
	)

	// Retrieve token details
	if token.TokenType == Email_Verification {
		existingToken, err = t.vtokenRepo.GetVCode(ctx, token.Email)
		if err != nil {
		return "", domain.ErrIncorrectEmail
	  	}

		if existingToken.Token != token.Token{
			return "", ErrIncorrectToken
		}
	    
	} else {
		existingToken, err = t.vtokenRepo.GetByToken(ctx, token.Token)
		if err != nil {
	     	return "", ErrIncorrectToken
	   }
	}
	
	if err != nil {
		return "", ErrIncorrectToken
	}

	if time.Now().After(existingToken.ExpiresAt) {
		return "", ErrExpiredToken
	}

	return existingToken.Email, nil
}

func (t *tokenUsecase) DeleteVCode(ctx context.Context, userID string) error {
	return t.vtokenRepo.DeleteVCode(ctx, userID)
}

func (t *tokenUsecase) FindByUserID(ctx context.Context, userID string) (*domain.Token, error) {
	return t.tokenRepo.FindByUserID(ctx, userID)
}

func (t *tokenUsecase) RefreshTokens(ctx context.Context, refreshToken string) (*domain.Token, error) {
	return t.tokenService.RefreshTokens(ctx, refreshToken)
}

func (t *tokenUsecase) GenerateTokens(ctx context.Context, userID string) (*domain.Token, error) {
	return t.tokenService.GenerateTokens(ctx, userID)
}

func (t *tokenUsecase) VerifyAccessToken(tokenString string) (string, error) {
	return t.tokenService.VerifyAccessToken(tokenString)
}

func (t *tokenUsecase) SaveToken(ctx context.Context, tokens *domain.Token) error{
	 return t.tokenRepo.Save(ctx, tokens)
}

func (t *tokenUsecase) DeleteByUserID(ctx context.Context, userID string) error {
	return t.tokenRepo.DeleteByUserID(ctx, userID)
}

func (t *tokenUsecase) GetByAccessToken(ctx context.Context, accessToken string)(string, error){
	return t.tokenRepo.FindByAccessToken(ctx, accessToken)


}
