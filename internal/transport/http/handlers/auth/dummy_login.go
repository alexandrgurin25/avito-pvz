package auth

import (
	"avito-pvz/internal/constants"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/transport/http/dto/auth"
	"encoding/json"
	"net/http"

	message "avito-pvz/internal/transport/http/dto/error"
)

func (h *authHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
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

	token, err := h.service.CreateDummyLogin(req.Role)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(auth.DummyLoginResponse{Token: token})
}
