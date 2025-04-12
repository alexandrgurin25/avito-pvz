package auth

type ErrorResponse struct {
	Message string `json:"message"`
}

type AuthHandler struct {
}

func NewHandler() *AuthHandler {
	return &AuthHandler{}
}
