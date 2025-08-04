package usecases

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
	"errors"

	domain "github.com/gedyzed/blog-starter-project/Domain"
)


type ITokenUsecase interface{
	CreateSendVCode(ctx context.Context, userID string, tokenType string)error
	GenerateSecureCode()(string, error)
	VerifyCode(ctx context.Context, userID string, vcode string) error
	DeleteVCode(ctx context.Context, userID string) error
	FindByUserID(ctx context.Context, userID string) (*domain.Token, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*domain.Token, error)
	GenerateTokens(ctx context.Context, userID string) (*domain.Token, error)
	VerifyAccessToken(string) (string, error)
	DeleteByUserID(ctx context.Context, userID string) error

}

type tokenUsecase struct {
	repo domain.ITokenRepo
	codeServices domain.ITokenService
	jwtServices domain.IJWTService
}

func NewTokenUsecase(repo domain.ITokenRepo, svs domain.ITokenService, js domain.IJWTService) ITokenUsecase {
	return &tokenUsecase{repo: repo, codeServices: svs, jwtServices:js }
} 

func(t *tokenUsecase) CreateSendVCode(ctx context.Context, userID string, tokenType string) error {

	// generate random verfication code 
	verfication_code, err := t.GenerateSecureCode()
	if err != nil {
		return err
	}

	// ten minutes of expiration time
	expiration_time := time.Now().Add(10 * time.Minute)

	token := domain.VToken {
				UserID: userID, 
				TokenType: tokenType, 
				Token: verfication_code,
				ExpiresAt: expiration_time,
			}
	
	// save the created verification code to db
	err = t.repo.CreateVCode(ctx, &token)
	if err != nil {
		fmt.Print("create", err.Error())
		return errors.New("internal server error")
	}

	// send the generate code via email
	subject := "Sending Verification Code"
	body := "Here is you verification code: " + verfication_code
	return t.codeServices.SendEmail([]string{userID,}, subject, body)
}

func (t *tokenUsecase) GenerateSecureCode()(string, error){

	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func (t *tokenUsecase) VerifyCode(ctx context.Context, userID string, vcode string)error {

	// retreive token details
	token, err := t.repo.GetVCode(ctx, userID)
	if err != nil {
		return errors.New("incorrect email or token")
	}

	if token.Token != vcode || time.Now().After(token.ExpiresAt) {
		return errors.New("invalid or expired verification code")
	}

	return nil
}

func (t *tokenUsecase) DeleteVCode(ctx context.Context, userID string) error {
	return t.repo.DeleteVCode(ctx, userID)
}

func (t *tokenUsecase) FindByUserID(ctx context.Context, userID string) (*domain.Token, error){
		return t.repo.FindByUserID(ctx, userID)
}

func (t *tokenUsecase) RefreshTokens(ctx context.Context, refreshToken string) (*domain.Token, error){
		return t.jwtServices.RefreshTokens(ctx, refreshToken)
} 

func (t *tokenUsecase) GenerateTokens(ctx context.Context, userID string) (*domain.Token, error){
		return t.jwtServices.GenerateTokens(ctx, userID)
}

func (t *tokenUsecase) VerifyAccessToken(tokenString string) (string, error){
	return t.jwtServices.VerifyAccessToken(tokenString)	
}

func (t *tokenUsecase) DeleteByUserID(ctx context.Context, userID string) error{
		return t.repo.DeleteByUserID(ctx, userID)
}


