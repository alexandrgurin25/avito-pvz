package reseption

import (
	myerrors "avito-pvz/internal/constants/errors"
	message "avito-pvz/internal/transport/http/dto/error"
	reception "avito-pvz/internal/transport/http/dto/reception"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *ReceptionHandler) CloseLastReception(w http.ResponseWriter, r *http.Request) {

	// Получаем параметр из URL
	pvzID := chi.URLParam(r, "pvzId")

	if pvzID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message.ErrorResponse{Message: myerrors.ErrPvzIdNil.Error()})
		return
	}

	closedReception, err := h.service.CloseLastReception(r.Context(), pvzID)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrActiveReceptionNotFound):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(message.ErrorResponse{Message: err.Error()})
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(message.ErrorResponse{Message: err.Error()})
			return
		}
	}

	response := reception.ReceptionResponse{
		ID:       closedReception.ID,
		PvzID:    closedReception.PvzID,
		DateTime: closedReception.DateTime.String(),
		Status:   closedReception.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
