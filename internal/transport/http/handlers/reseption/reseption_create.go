package reseption

import (
	myerrors "avito-pvz/internal/constants/errors"
	message "avito-pvz/internal/transport/http/dto/error"
	"avito-pvz/internal/transport/http/dto/reception"
	"avito-pvz/pkg/logger"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (h *ReceptionHandler) CreateReception(w http.ResponseWriter, r *http.Request) {
	var req reception.CreateReceptionRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx,
			"Failed to decode JSON request",
			zap.Error(err),
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.Any("headers", r.Header),
		)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrPvzIdNil.Error()})
		return
	}

	if req.PvzID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrPvzIdNil.Error()})
		return
	}

	createdReception, err := h.service.CreateReception(ctx, req.PvzID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: err.Error()})
		return
	}

	res := &reception.ReceptionResponse{
		ID:       createdReception.ID,
		PvzID:    createdReception.PvzID,
		DateTime: createdReception.DateTime.Format(time.RFC3339),
		Status:   createdReception.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "Failed to encode response:", zap.Error(err))
	}
}
