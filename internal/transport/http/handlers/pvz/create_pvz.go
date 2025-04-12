package pvz

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"avito-pvz/internal/transport/http/dto/pvz"
	"avito-pvz/pkg/logger"
	"encoding/json"
	"errors"
	"log"
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
		json.NewEncoder(w).Encode(ErrorResponse{Message: myerrors.ErrCityOrIdNil.Error()})
		return
	}

	if req.ID == "" || req.City == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: myerrors.ErrCityOrIdNil.Error()})
		return
	}

	newPVZ := &entity.PVZ{
		UUID: req.ID,
		City: entity.City{
			Name: req.City,
		},
	}

	pvz, err := h.service.CreatePVZ(ctx, newPVZ)

	if err != nil {
		if errors.Is(err, myerrors.ErrCityNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: myerrors.ErrCityNotFound.Error()})
			return
		}

		if errors.Is(err, myerrors.ErrCityNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: myerrors.ErrCityNotFound.Error()})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		return
	}

	res.City = pvz.City.Name
	res.ID = pvz.UUID
	res.RegistrationDate = pvz.CreatedAt.Format("2006-01-02T15:04:05.000Z")

	w.WriteHeader(http.StatusCreated) // 201
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
