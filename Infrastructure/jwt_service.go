package infrastructure

import (
	"fmt"
	"time"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	key string
}

func (jwts *JWTService) GenerateToken() (*domain.Token, error) {
	t1 := jwt.New(jwt.SigningMethodES256)
	t2 := jwt.New(jwt.SigningMethodES256)

	accessToken, err := t1.SignedString([]byte(jwts.key))
	if err != nil {
		return nil, err
	}

	refreshToken, err := t2.SignedString([]byte(jwts.key))
	if err != nil {
		return nil, err
	}

	token := domain.Token{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(30 * (24 * time.Hour)),
	}

	return &token, nil
}

func (jwts *JWTService) ValidateToken(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwts.key), nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (jwts *JWTService) RefreshToken(token string) (string, error) {
	t1 := jwt.New(jwt.SigningMethodES256)

	accessToken, err := t1.SignedString([]byte(jwts.key))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
