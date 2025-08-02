package infrastructure

import (
	domain "github.com/gedyzed/blog-starter-project/Domain"
	"golang.org/x/crypto/bcrypt"
)

type passwordService struct{}

func (ps *passwordService) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (ps *passwordService) Verify(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hashedPassword))
}

func NewPasswordService() domain.IPasswordService {
	return &passwordService{}
}
