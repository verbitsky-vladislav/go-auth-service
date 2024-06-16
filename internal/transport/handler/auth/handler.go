package auth

import "auth-microservice/internal/service"

type Handler struct {
	userService service.UserService
}

func NewAuthHandler(
	userService service.UserService,
) *Handler {
	return &Handler{
		userService: userService,
	}
}
