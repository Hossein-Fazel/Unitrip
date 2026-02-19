package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID, role string) (string, error)
	ValidateToken(token string) (string, string, error) // userID, role
}

type jwtService struct {
	secretKey string
	ttl       time.Duration
}

func NewJWTService(secretKey string, ttl time.Duration) JWTService {
	return &jwtService{
		secretKey: secretKey,
		ttl:       ttl,
	}
}

func (j *jwtService) GenerateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(j.ttl).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) ValidateToken(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", "", errors.New("user_id not found in token")
		}
		role, ok := claims["role"].(string)
		if !ok {
			return "", "", errors.New("role not found in token")
		}
		return userID, role, nil
	}

	return "", "", errors.New("invalid token")
}
