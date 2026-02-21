package jwt

import (
	"errors"
	"time"

	"unitrip/internal/adapter/http"
	"unitrip/internal/entity"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey string
	ttl       time.Duration
}

func NewJWTService(secretKey string, ttl time.Duration) http.JWTService {
	return &jwtService{
		secretKey: secretKey,
		ttl:       ttl,
	}
}

func (j *jwtService) GenerateToken(user *entity.User) (string, error) {
	claims := http.AuthClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) ValidateToken(tokenStr string) (*http.AuthClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&http.AuthClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(j.secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*http.AuthClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
