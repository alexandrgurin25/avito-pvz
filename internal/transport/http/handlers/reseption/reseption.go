package reseption

import "avito-pvz/internal/service/reception"

type ReceptionHandler struct {
	service reception.Service
}

func NewHandler(service reception.Service) *ReceptionHandler {
	return &ReceptionHandler{service: service}
}
