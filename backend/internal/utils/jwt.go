package utils

import (
	"time"

	"streaming-system/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string          `json:"userId"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(secret string, userID string, role models.UserRole, duration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
