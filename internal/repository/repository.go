package repository

import (
	"auth-microservice/internal/model"
)

type UserRepository interface {
	Create(user *model.UserCreate) (string, error)
	CreateFromGoogle(user *model.UserCreateFromGoogle) (string, error)
	Update(id string, user *model.UserUpdate) error
	FindById(id string) (*model.User, error)
	FindByEmail(id string) (*model.User, error)
}
