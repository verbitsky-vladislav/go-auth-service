package base_auth_responses

type CreateUserResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}
