package infrastructure

type ITokenService interface {
	GenerateToken() (string, error)
	ValidateToken(string) error
}

type JWTTokenSevice struct {
}
