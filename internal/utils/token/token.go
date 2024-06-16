package token

import (
	"auth-microservice/pkg/logger"
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomToken генерирует случайный токен длиной 32 символа.
func GenerateRandomToken() (string, error) {
	tokenLength := 32

	randomBytes := make([]byte, tokenLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", logger.Error(err, "failed to generate random token")
	}
	token := base64.URLEncoding.EncodeToString(randomBytes)
	token = token[:tokenLength]

	return token, nil
}
