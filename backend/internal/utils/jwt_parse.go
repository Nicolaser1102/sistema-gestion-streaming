package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(secret string, tokenString string) (*JWTClaims, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validar algoritmo
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("algoritmo inválido")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(*JWTClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("token inválido")
	}
	return claims, nil
}
