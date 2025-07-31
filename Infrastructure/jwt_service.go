package infrastructure

import (
	"github.com/gedyzed/blog-starter-project/Domain"
)

type ITokenService interface {
	GenerateToken() (domain.Token, error)
	ValidateToken(string) error
	RefreshToken(string) (string, error)
}

type JWTTokenSevice struct {
}
