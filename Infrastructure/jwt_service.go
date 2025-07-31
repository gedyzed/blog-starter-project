package infrastructure

type ITokenService interface {
	GenerateToken() (string, error)
	ValidateToken(string) error
	RefreshToken(string) (string, error)
}

type JWTTokenSevice struct {
}
