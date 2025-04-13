package pvz

import (
	message "avito-pvz/internal/transport/http/dto/error"
	"avito-pvz/internal/transport/http/dto/pvz"
	"avito-pvz/pkg/logger"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (h *PvzHandler) GetPVZs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req pvz.GetPVZsRequest
	if err := req.ParseFromQuery(r.URL.Query()); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to parse query params",
			zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{
			Message: "invalid query parameters",
		})
		return
	}

	// Валидация параметров
	if req.Limit < 1 || req.Limit > 30 {
		req.Limit = 10
	}
	if req.Page < 1 {
		req.Page = 1
	}

	pvzs, err := h.service.GetAllWithReceptions(ctx, req.StartDate, req.EndDate, req.Page, req.Limit)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get pvzs",
			zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := make([]pvz.PVZWithReceptions, 0, len(pvzs))
	for _, p := range pvzs {
		response = append(response, pvz.PVZWithReceptions{
			PVZ: pvz.PvzResponse{
				ID:               p.UUID,
				RegistrationDate: p.CreatedAt,
				City:             p.City.Name,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to encode response",
			zap.Error(err))
	}
}
