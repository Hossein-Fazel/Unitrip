package http

import (
	"unitrip/internal/entity"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(user *entity.User) (string, error)
	ValidateToken(token string) (*AuthClaims, error)
}
type AuthClaims struct {
	UserID int64           `json:"user_id"`
	Role   entity.UserRole `json:"role"`
	jwt.RegisteredClaims
}
