package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plaintext password using bcrypt
func HashPassword(plaintext string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CompareHashedPassword compares a bcrypt-hashed password with a plaintext one
func CompareHashedPassword(hashedPass string, plainPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainPass)) == nil
}
