package auth

import (
	"auth-microservice/internal/config"
	"auth-microservice/internal/service"
)

type Handler struct {
	cfg         *config.Config
	userService service.UserService
	authService service.AuthService
}

func NewAuthHandler(
	cfg *config.Config,
	userService service.UserService,
	authService service.AuthService,
) *Handler {
	return &Handler{
		cfg:         cfg,
		userService: userService,
		authService: authService,
	}
}
