package infrastructure

import (
	domain "github.com/gedyzed/blog-starter-project/Domain"
	"golang.org/x/crypto/bcrypt"
)

type passwordService struct{}

func NewPasswordService() domain.IPasswordService {
	return &passwordService{}

}

func (ps *passwordService) Verify(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Hash password
func (s *passwordService) Hash(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	password = string(hashedPassword)
	return password, nil
}
