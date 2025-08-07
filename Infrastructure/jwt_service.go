package infrastructure

import (
	"context"
	"time"

	"github.com/gedyzed/blog-starter-project/Domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenService struct {
	repo       domain.ITokenRepo
	accessKey  []byte
	refreshKey []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTTokenService(repo domain.ITokenRepo, accessKey, refreshKey string, accessTTL, refreshTTL time.Duration) *JWTTokenService {
	return &JWTTokenService{
		repo,
		[]byte(accessKey),
		[]byte(refreshKey),
		accessTTL,
		refreshTTL,
	}
}

func (s *JWTTokenService) signJWT(userID string, key []byte, ttl time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func (s *JWTTokenService) verifyJWT(tokenString string, key []byte) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return key, nil
	})

	if err != nil {
		return "", domain.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", domain.ErrInvalidToken
}

func (s *JWTTokenService) GenerateTokens(ctx context.Context, userID string) (*domain.Token, error) {
	accessToken, err := s.signJWT(userID, s.accessKey, s.accessTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.signJWT(userID, s.refreshKey, s.refreshTTL)
	if err != nil {
		return nil, err
	}

	tokens := domain.Token{
		UserID:        userID,
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpiry:  time.Now().Add(s.accessTTL),
		RefreshExpiry: time.Now().Add(s.refreshTTL),
	}

	if err := s.repo.Save(ctx, &tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}

func (s *JWTTokenService) RefreshTokens(ctx context.Context, refreshToken string) (*domain.Token, error) {
	userID, err := s.verifyJWT(refreshToken, s.refreshKey)
	if err != nil {
		return nil, domain.ErrInvalidRefreshToken
	}

	tokens, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, domain.ErrInvalidRefreshToken
	}

	if tokens.RefreshToken != refreshToken {
		return nil, domain.ErrInvalidRefreshToken
	}

	if tokens.RefreshExpiry.Before(time.Now()) {
		_ = s.repo.DeleteByUserID(ctx, userID)
		return nil, domain.ErrInvalidRefreshToken
	}

	return s.GenerateTokens(ctx, userID)
}

func (s *JWTTokenService) VerifyAccessToken(tokenString string) (string, error) {
	return s.verifyJWT(tokenString, s.accessKey)
}
