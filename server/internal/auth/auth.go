package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GenerateRefreshToken() (string, error)
}
