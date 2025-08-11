package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"strings"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	usecases "github.com/gedyzed/blog-starter-project/Usecases"
	"golang.org/x/oauth2"
)

type OAuthServices struct {
	config 		*oauth2.Config
	userUsecase *usecases.UserUsecases
}

func NewOAuthServices(cfg *oauth2.Config,  uc *usecases.UserUsecases) domain.IOAuthServices {
	return &OAuthServices{config:cfg, userUsecase: uc}
}

func (os OAuthServices) VerifyGoogleIDToken(ctx context.Context, accessToken string)(string, error){
	
	userID, err := os.userUsecase.GetToken(ctx, accessToken)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (os OAuthServices) RefreshToken(ctx context.Context, token *domain.Token)(*domain.Token, error){

	expiredToken := &oauth2.Token{
		RefreshToken: token.RefreshToken,
		Expiry:       time.Now().Add(-time.Hour),
	}

	tokenSource := os.config.TokenSource(ctx, expiredToken)

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, domain.ErrFailedToFetch
	}

	token.AccessToken = newToken.AccessToken
	token.AccessExpiry = newToken.Expiry
	token.UpdatedAt = time.Now()

	if newToken.RefreshToken != "" {
		token.RefreshToken = newToken.RefreshToken
	}

	err = os.userUsecase.SaveToken(ctx, token)
	if err != nil {
		return nil, err
	}
	
	return token, nil
}

func (os OAuthServices) ResolveUserID(ctx context.Context, email string)(string, error){

	existing, err := os.userUsecase.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	userID := existing.ID.Hex()
	return userID, nil
}

func (os OAuthServices) OAuthCallBack(ctx context.Context, code string) (*domain.Token, error){

	got, err := os.config.Exchange(ctx, code)
	if err != nil {
		return nil, domain.ErrFailedToExchange
	}

	refreshTTL := 30 * 24 * time.Hour // 30 days of refresh token expiry
	now := time.Now()

	tokens := &domain.Token{
		AccessToken:   got.AccessToken,
		RefreshToken:  got.RefreshToken,
		AccessExpiry:  got.Expiry,
		RefreshExpiry: now.Add(refreshTTL),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	client := os.config.Client(ctx, got)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, domain.ErrFailedToFetchUserInfo
	}
	defer resp.Body.Close() 

	var userInfo domain.UserInfo // Use value instead of pointer
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, domain.ErrFailedToDecodeUserInfo
	}

	existingUser, err := os.userUsecase.GetByEmail(ctx, userInfo.Email)
	var userID string

	if err == nil && existingUser != nil {
		userID = existingUser.ID.Hex()
	} else if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	} else {

		parts := strings.Split(userInfo.Email, "@")
		username := parts[0]

		profile := domain.Profile{
			ProfilePic: userInfo.Picture,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		user := &domain.User{
			Firstname: userInfo.GivenName,
			Lastname:  userInfo.FamilyName,
			Username: username,
			Email:     userInfo.Email,
			Role:      "user",
			Provider:  "google",
			CreatedAt: now,
			UpdatedAt: now,

			Profile: profile,
		}

		userID, err = os.userUsecase.Register(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	tokens.UserID = userID
	err = os.userUsecase.SaveToken(ctx, tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}





