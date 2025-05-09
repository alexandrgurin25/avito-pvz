package pvz

import (
	pvzS "avito-pvz/internal/service/pvz"
)

type PvzHandler struct {
	service pvzS.PVZService
}

func New(service pvzS.PVZService) *PvzHandler {
	return &PvzHandler{service: service}
}
