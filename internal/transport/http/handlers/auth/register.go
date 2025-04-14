package auth

import (
	"avito-pvz/internal/constants"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/transport/http/dto/auth"
	message "avito-pvz/internal/transport/http/dto/error"
	"avito-pvz/pkg/logger"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req auth.RegisterRequest

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

	if req.Email == "" || req.Password == "" || req.Role == "" {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Username, password and role is required")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrEmailOrPasswordEmpty.Error()})
		return
	}
	
	if req.Role != constants.Employee && req.Role != constants.Moderator {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrInvalidRole.Error()})
		return
	}

	err := h.service.Register(ctx, req.Email, req.Password, req.Role)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
