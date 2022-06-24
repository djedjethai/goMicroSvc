package dto

type JsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type RequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
