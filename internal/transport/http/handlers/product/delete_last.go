package product

import (
	myerrors "avito-pvz/internal/constants/errors"
	message "avito-pvz/internal/transport/http/dto/error"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *ProductHandler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	pvzID := chi.URLParam(r, "pvzId")
	ctx := r.Context()

	if pvzID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrPvzIdNil.Error()})
		return
	}

	err := h.service.DeleteLastProduct(ctx, pvzID)

	if err != nil {
		if errors.Is(err, myerrors.ErrPVZNotFound) || errors.Is(err, myerrors.ErrActiveReceptionNotFound) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(http.StatusOK)
}
