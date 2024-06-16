package service

import (
	"auth-microservice/internal/model"
	"time"
)

type UserService interface {
	CreateUser(user *model.UserCreate) (string, error)
	UpdateUser(id string, user *model.UserUpdate) error
	FindUserById(id string) (*model.User, error)
	FindUserByEmail(id string) (*model.User, error)
}

// MailerService added from pkg/mailer
type MailerService interface {
	SendMail(dest []string, subject, bodyMessage string) error
	SendVerificationUrl(dest []string, username, verificationUrl string) error
	WriteEmail(dest []string, contentType, subject, bodyMessage string) (string, error)
	WriteHTMLEmail(dest []string, subject, bodyMessage string) (string, error)
	WritePlainEmail(dest []string, subject, bodyMessage string) (string, error)
}

type CacheService interface {
	SetExpire(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type AuthService interface {
	Register(user *model.UserCreate) (string, error)
	Login() error
	Logout() error
}
