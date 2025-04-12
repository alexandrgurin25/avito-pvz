package product

import "avito-pvz/internal/service/product"

type ProductHandler struct {
	service product.Service
}

func New(service product.Service) *ProductHandler {
	return &ProductHandler{service: service}
}
