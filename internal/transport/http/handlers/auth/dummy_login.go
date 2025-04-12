package auth

import (
	"avito-pvz/internal/constants"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/transport/http/dto/auth"
	"encoding/json"
	"net/http"
	"time"

	message "avito-pvz/internal/transport/http/dto/error"

	"github.com/golang-jwt/jwt"
)

func (h *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req auth.DummyLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: err.Error()})
		return
	}

	if req.Role != constants.Employee && req.Role != constants.Moderator {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrInvalidRole.Error()})
		return
	}

	id := 1 //Фиктивный ID
	token, err := generateJWT(id, &req)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrInvalidGenerateJWT.Error()})
		return
	}

	json.NewEncoder(w).Encode(auth.DummyLoginResponse{Token: token})
}

func generateJWT(id int, req *auth.DummyLoginRequest) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"role":    &req.Role,
		"iat":     time.Now().Add(constants.ExpirationTime).Unix(),
	})

	tokenString, err := token.SignedString([]byte("avito"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
