package service

import (
	"auth-microservice/internal/model"
	"golang.org/x/oauth2"
	"time"
)

type UserService interface {
	CreateUser(user *model.UserCreate) (string, error)
	CreateUserFromGoogle(user *model.UserCreateFromGoogle) (string, error)
	UpdateUser(id string, user *model.UserUpdate) error
	FindUserById(id string) (*model.User, error)
	FindUserByEmail(id string) (*model.User, error)
}

type MailerService interface {
	SendMail(dest []string, subject, bodyMessage string) error
	SendVerificationUrl(dest []string, username, verificationUrl string) error
	WriteEmail(dest []string, contentType, subject, bodyMessage string) (string, error)
	WriteHTMLEmail(dest []string, subject, bodyMessage string) (string, error)
	WritePlainEmail(dest []string, subject, bodyMessage string) (string, error)
}

type JwtService interface {
	GenerateTokens(user model.UserInfo) (model.Tokens, error)
	VerifyAccessToken(tokenString string) (model.UserInfo, error)
	VerifyRefreshToken(tokenString string) error
	GenerateRefreshToken() (string, error)
	RefreshTokens(refreshToken string) (model.Tokens, error)
}

type CacheService interface {
	SetExpire(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type AuthService interface {
	Register(user *model.UserCreate) (string, error)
	Login(user *model.UserLogin) (*model.UserInfo, *model.Tokens, error)
	Logout() error

	ConfirmEmail(token string) (*model.UserInfo, error)
	SendVerificationEmail(info *model.UserInfo) error
}

type GoogleService interface {
	GetGoogleConfig() *oauth2.Config
}
