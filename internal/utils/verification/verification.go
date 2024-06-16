package verification

import (
	"auth-microservice/internal/utils/token"
	"auth-microservice/pkg/logger"
)

func GetVerificationUrl(appUrl string) (string, string, error) {
	verificationToken, err := token.GenerateRandomToken()
	if err != nil {
		return "", "", logger.Error(err, "error generating token")
	}
	return appUrl + "/confirm/verify/" + verificationToken, verificationToken, nil
}
