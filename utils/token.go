package utils

import (
	"errors"
	"personal-finance/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenParams holds parameters needed to generate a token
type TokenParams struct {
	ID string
	// Role string // e.g. "admin" | "user"
}

// CustomClaims is the JWT claims structure
type CustomClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

// ParseToken generates a JWT token with 24h expiry
func ParseToken(params TokenParams) (string, error) {
	claims := jwt.MapClaims{
		"id":  params.ID,
		"sub": params.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		// "role": params.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Application.JWT.SecretKey))
}

// ParseTokenWithExpiration generates a JWT token with a custom expiry duration
func ParseTokenWithExpiration(params TokenParams, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":  params.ID,
		"sub": params.ID,
		"exp": time.Now().Add(expiration).Unix(),
		// "role": params.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Application.JWT.SecretKey))
}

// VerifyToken decodes and validates a JWT token string
func VerifyToken(tokenStr string) (*CustomClaims, error) {
	secret := config.Application.JWT.SecretKey
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token sudah kadaluarsa, silahkan melakukan login ulang")
	}

	return claims, nil
}
