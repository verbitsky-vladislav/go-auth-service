package auth

import (
	"auth-microservice/internal/model"
	"auth-microservice/internal/utils/constant"
	"auth-microservice/pkg/logger"
)

func (a authService) ConfirmEmail(token string) error {
	userInfo, err := a.verificationService.Decrypt(token)
	if err != nil {
		return err
	}

	cacheToken, err := a.cacheService.Get(userInfo.ID)
	if err != nil {
		return err
	}
	result := a.verificationService.CompareTokens(token, cacheToken)
	if !result {
		return logger.Error(nil, "tokens don't match")
	}

	err = a.userService.UpdateUser(userInfo.ID, &model.UserUpdate{
		IsVerified: result,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a authService) SendVerificationEmail(info *model.UserInfo) error {
	url, token, err := a.verificationService.GetVerificationUrl(*info)
	if err != nil {
		return err
	}

	err = a.cacheService.SetExpire(info.ID, token, constant.OTPTokenLifeDuration)
	if err != nil {
		return err
	}

	err = a.mailerService.SendVerificationUrl(
		[]string{info.Email},
		info.Username,
		url,
	)
	if err != nil {
		return logger.Error(err, "send verification url failed")
	}
	return nil
}
