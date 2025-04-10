package auth

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Role     string `json:"role"`
}

type RegisterResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Role  string `json:"role"`
}