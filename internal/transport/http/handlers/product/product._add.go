package product

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/pkg/logger"
	"encoding/json"
	"net/http"

	message "avito-pvz/internal/transport/http/dto/error"
	"avito-pvz/internal/transport/http/dto/product"

	"go.uber.org/zap"
)

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var req *product.AddProductRequest

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
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrJsonNotFound.Error()})
		return
	}

	if req.PvzId == "" || req.Type == ""{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrPvzIdOrTypeNil.Error()})
		return
	}

	createdProduct, err := h.service.AddProduct(ctx, req.Type, req.PvzId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: err.Error()})
		return
	}

	res := &product.ProductResponse{
		ID: createdProduct.ID,
		Type: createdProduct.Category,
		DateTime: createdProduct.DateTime,
		ReceptionID: createdProduct.ReceptionID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "Failed to encode response:", zap.Error(err))
	}
}
