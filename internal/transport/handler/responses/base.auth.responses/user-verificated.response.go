package base_auth_responses

import "auth-microservice/internal/model"

type VerificatedUserResponse struct {
	Status   int             `json:"status"`
	Message  string          `json:"message"`
	UserInfo *model.UserInfo `json:"user_info"`
}
