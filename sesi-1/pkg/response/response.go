package response

type ApiResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    interface{}       `json:"data"`
	Errors  map[string]string `json:"errors,omitempty"`
}