package pvz

import (
	message "avito-pvz/internal/transport/http/dto/error"
	"avito-pvz/internal/transport/http/dto/product"
	"avito-pvz/internal/transport/http/dto/pvz"
	"avito-pvz/internal/transport/http/dto/reception"
	"avito-pvz/pkg/logger"
	"encoding/json"
	"net/http"
	"time"

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
		receptions := make([]pvz.ReceptionWithProducts, 0, len(p.Receptions))
		for _, rec := range p.Receptions {
			products := make([]product.ProductResponse, 0, len(rec.Products))
			for _, prod := range rec.Products {
				products = append(products, product.ProductResponse{
					ID:          prod.ID,
					ReceptionID: prod.ReceptionID,
					DateTime:    prod.DateTime,
					Type:        prod.Category,
				})
			}

			receptions = append(receptions, pvz.ReceptionWithProducts{
				Reception: reception.ReceptionResponse{
					ID:       rec.ID,
					PvzID:    rec.PvzID,
					DateTime: rec.DateTime.Truncate(time.Millisecond).Format("2006-01-02T15:04:05.000Z"),
					Status:   rec.Status,
				},
				Products: products,
			})
		}

		response = append(response, pvz.PVZWithReceptions{
			PVZ: pvz.PvzResponse{
				ID:               p.UUID,
				RegistrationDate: p.CreatedAt,
				City:             p.City.Name,
			},
			Receptions: receptions,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to encode response",
			zap.Error(err))
	}
}
