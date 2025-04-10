package auth

import (
	errors "avito-pvz/internal/transport/http/dto"
	"avito-pvz/internal/transport/http/dto/auth"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

const employee = "employee"
const moderator = "moderator"

type ErrorResponse struct {
	Message string `json:"message"`
}

func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req auth.DummyLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		return
	}

	if req.Role != employee && req.Role != moderator {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: errors.ErrInvalidRole.Error()})
		return
	}

	id := 1 //Фиктивный ID
	token, err := generateJWT(id, &req)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: errors.ErrInvalidGenerateJWT.Error()})
		return
	}

	json.NewEncoder(w).Encode(auth.DummyLoginResponse{Token: token})
}

func generateJWT(id int, req *auth.DummyLoginRequest) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"role":    &req.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("avito_test_secret"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
