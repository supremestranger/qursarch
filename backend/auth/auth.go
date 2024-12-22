// auth/auth.go
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword хэширует пароль с использованием SHA-256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// GenerateSessionID генерирует уникальный идентификатор сессии
func GenerateSessionID() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
