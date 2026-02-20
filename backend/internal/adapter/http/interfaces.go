package http

type JWTService interface {
	GenerateToken(userID, role string) (string, error)
	ValidateToken(token string) (string, string, error) // userID, role
}