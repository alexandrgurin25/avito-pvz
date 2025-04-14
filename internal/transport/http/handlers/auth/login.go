package auth

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/transport/http/dto/auth"
	message "avito-pvz/internal/transport/http/dto/error"
	"avito-pvz/pkg/logger"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req auth.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx,
			"Failed to decode JSON request",
			zap.Error(err),
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.Any("headers", r.Header),
		)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrJsonNotFound.Error()})
		return
	}

	if req.Email == "" || req.Password == "" {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Username and password is required")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrEmailOrPasswordEmpty.Error()})
		return
	}

	token, err := h.service.Login(ctx, req.Email, req.Password)

	var res auth.LoginResponse

	res.Token = token

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: err.Error()})
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "error encoding JSON", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
