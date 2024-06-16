package errors

type ErrorResponse struct {
	Error Error `json:"error"`
}
type Error struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Error     string `json:"error"`
	Message   string `json:"message"`
}
