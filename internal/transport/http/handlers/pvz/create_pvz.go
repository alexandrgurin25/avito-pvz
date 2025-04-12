package pvz

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	message "avito-pvz/internal/transport/http/dto/error"
	"avito-pvz/internal/transport/http/dto/pvz"
	"avito-pvz/pkg/logger"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

func (h *PvzHandler) CreatePVZ(w http.ResponseWriter, r *http.Request) {
	var req pvz.PvzRequest
	var res pvz.PvzResponse

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
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrCityOrIdNil.Error()})
		return
	}

	if req.ID == "" || req.City == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrCityOrIdNil.Error()})
		return
	}

	newPVZ := &entity.PVZ{
		UUID: req.ID,
		City: entity.City{
			Name: req.City,
		},
		CreatedAt: req.RegistrationDate,
	}

	pvz, err := h.service.CreatePVZ(ctx, newPVZ)

	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, myerrors.ErrCityNotFound) {
			status = http.StatusBadRequest
		}

		if errors.Is(err, myerrors.ErrCityNotFound) {
			status = http.StatusBadRequest
		}

		if errors.Is(err, myerrors.ErrPVZAlreadyExists) {
			status = http.StatusBadRequest
		}

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: err.Error()})
		return
	}

	res.City = pvz.City.Name
	res.ID = pvz.UUID
	res.RegistrationDate = pvz.CreatedAt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "Failed to encode response:", zap.Error(err))
	}
}
