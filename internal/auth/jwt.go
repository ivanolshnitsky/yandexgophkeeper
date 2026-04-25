package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService сервис работы с JWT.
type JWTService struct {
	secret string
	ttl    time.Duration
}

// NewJWTService создаёт JWTService.
func NewJWTService(secret string, ttl time.Duration) *JWTService {
	return &JWTService{
		secret: secret,
		ttl:    ttl,
	}
}

// Generate создаёт JWT-токен.
func (j *JWTService) Generate(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(j.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}
