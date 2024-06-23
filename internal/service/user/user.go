package user

import (
	"auth-microservice/internal/model"
	"auth-microservice/internal/repository"
	"auth-microservice/internal/service"
	"auth-microservice/pkg/logger"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) service.UserService {
	return &userService{userRepo: userRepo}
}

func (u userService) CreateUser(user *model.UserCreate) (string, error) {
	var checkUser *model.User
	checkUser, err := u.FindUserByEmail(user.Email)
	if err != nil {
		return "", err
	}
	if checkUser != nil {
		return "", logger.Error(err, "user already exists")
	}

	id, err := u.userRepo.Create(user)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (u userService) CreateUserFromGoogle(user *model.UserCreateFromGoogle) (string, error) {
	var checkUser *model.User
	checkUser, err := u.FindUserByEmail(user.Email)
	if err != nil {
		return "", err
	}
	if checkUser != nil {
		return "", logger.Error(err, "user already exists")
	}

	id, err := u.userRepo.CreateFromGoogle(user)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (u userService) UpdateUser(id string, user *model.UserUpdate) error {
	err := u.userRepo.Update(id, user)
	if err != nil {
		return err
	}
	return nil
}

func (u userService) FindUserById(id string) (*model.User, error) {
	user, err := u.userRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u userService) FindUserByEmail(id string) (*model.User, error) {
	user, err := u.userRepo.FindByEmail(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
