package auth

import (
	"auth-microservice/internal/config"
	"auth-microservice/internal/model"
	"auth-microservice/internal/service"
	"auth-microservice/internal/utils/verification"
	"auth-microservice/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	cfg                 *config.Config
	userService         service.UserService
	cacheService        service.CacheService
	mailerService       service.MailerService
	verificationService verification.VerificationService
}

func NewAuthService(
	userService service.UserService,
	cacheService service.CacheService,
	mailerService service.MailerService,
	verificationService verification.VerificationService,
) service.AuthService {
	return &authService{
		userService:         userService,
		cacheService:        cacheService,
		mailerService:       mailerService,
		verificationService: verificationService,
	}
}

func (a authService) Register(user *model.UserCreate) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", logger.Error(err, "failed to hash password")
	}
	user.Password = string(hashedPassword)

	id, err := a.userService.CreateUser(user)
	if err != nil {
		return "", err
	}

	err = a.SendVerificationEmail(&model.UserInfo{
		ID:       id,
		Email:    user.Email,
		Username: user.Username,
	})
	if err != nil {
		return "", err
	}

	return id, nil
}

func (a authService) Login(user *model.UserLogin) error {
	checkUser, err := a.userService.FindUserByEmail(user.Email)
	if err != nil {
		return err
	}
	if checkUser == nil {
		return logger.Error(nil, "user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(user.Password))
	if err != nil {
		return logger.Error(nil, "incorrect password")
	}
	if !checkUser.IsVerified {
		err = a.SendVerificationEmail(&model.UserInfo{
			ID:       checkUser.Id,
			Email:    checkUser.Email,
			Username: checkUser.Username,
		})
		if err != nil {
			return err
		}
		return logger.Error(nil, "user is not verified. verification url was sent on email")
	}

	//todo add jwt
	return err

}

func (a authService) Logout() error {
	//TODO implement me
	panic("implement me")
}
