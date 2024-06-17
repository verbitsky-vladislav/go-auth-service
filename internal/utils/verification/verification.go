package verification

import (
	"auth-microservice/internal/config"
	"auth-microservice/internal/model"
	"auth-microservice/pkg/logger"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

type VerificationService interface {
	Encrypt(info model.UserInfo) (string, error)
	Decrypt(encrypted string) (model.UserInfo, error)
	GetVerificationUrl(info model.UserInfo) (string, string, error)
	CompareTokens(token1, token2 string) bool
}

type verificationService struct {
	cfg *config.Config
}

func NewVerificationService(
	cfg *config.Config,
) VerificationService {
	return &verificationService{
		cfg: cfg,
	}
}

func (v verificationService) Encrypt(info model.UserInfo) (string, error) {
	data, err := json.Marshal(info)
	if err != nil {
		return "", logger.Error(err, "failed to marshal info")
	}

	block, err := aes.NewCipher([]byte(v.cfg.Application.VERIFICATION_PASSPHRASE))
	if err != nil {
		return "", logger.Error(err, "failed to create cipher | block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", logger.Error(err, "failed to create cipher | gcm")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", logger.Error(err, "failed to generate nonce")
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (v verificationService) Decrypt(encrypted string) (model.UserInfo, error) {
	data, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return model.UserInfo{}, logger.Error(err, "failed to decode encrypted data")
	}

	block, err := aes.NewCipher([]byte(v.cfg.Application.VERIFICATION_PASSPHRASE))
	if err != nil {
		return model.UserInfo{}, logger.Error(err, "failed to create cipher | block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return model.UserInfo{}, logger.Error(err, "failed to create cipher | gcm")
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return model.UserInfo{}, logger.Error(err, "nonce is too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return model.UserInfo{}, logger.Error(err, "failed to decrypt data")
	}

	var info model.UserInfo
	err = json.Unmarshal(plaintext, &info)
	if err != nil {
		return model.UserInfo{}, logger.Error(err, "failed to unmarshal data")
	}

	return info, nil
}

func (v verificationService) CompareTokens(token1, token2 string) bool {
	return token1 == token2
}

func (v verificationService) GetVerificationUrl(info model.UserInfo) (string, string, error) {
	verificationToken, err := v.Encrypt(info)
	if err != nil {
		return "", "", fmt.Errorf("error generating token: %w", err)
	}
	return v.cfg.Application.URL + "/confirm/verify/" + verificationToken, verificationToken, nil
}
