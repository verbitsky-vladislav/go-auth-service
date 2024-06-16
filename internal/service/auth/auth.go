package auth

import (
	"auth-microservice/internal/config"
	"auth-microservice/internal/model"
	"auth-microservice/internal/service"
	"auth-microservice/internal/utils/constant"
	"auth-microservice/internal/utils/verification"
	"auth-microservice/pkg/logger"
)

type authService struct {
	cfg           *config.Config
	userService   service.UserService
	cacheService  service.CacheService
	mailerService service.MailerService
}

func NewAuthService(
	userService service.UserService,
	cacheService service.CacheService,
	mailerService service.MailerService,
) service.AuthService {
	return &authService{
		userService:   userService,
		cacheService:  cacheService,
		mailerService: mailerService,
	}
}

func (a authService) Register(user *model.UserCreate) (string, error) {
	id, err := a.userService.CreateUser(user)
	if err != nil {
		return "", err
	}

	url, token, err := verification.GetVerificationUrl(a.cfg.Application.URL)
	if err != nil {
		return "", err
	}

	err = a.cacheService.SetExpire(token, id, constant.OTPTokenLifeDuration)
	if err != nil {
		return "", err
	}

	err = a.mailerService.SendVerificationUrl(
		[]string{user.Email},
		user.Username,
		url,
	)
	if err != nil {
		return "", logger.Error(err, "send verification url failed")
	}

	return id, nil
}

func (a authService) Login() error {
	//TODO implement me
	panic("implement me")
}

func (a authService) Logout() error {
	//TODO implement me
	panic("implement me")
}
