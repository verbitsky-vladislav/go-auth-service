package responses

type LogoutUserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
