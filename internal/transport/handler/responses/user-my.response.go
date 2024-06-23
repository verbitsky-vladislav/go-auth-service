package responses

import "auth-microservice/internal/model"

type UserMyResponse struct {
	Status   int         `json:"status"`
	Message  string      `json:"message"`
	UserInfo *model.User `json:"user"`
}
