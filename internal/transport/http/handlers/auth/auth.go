package auth

import "avito-pvz/internal/service/auth"

type authHandler struct {
	service auth.Service
}

func NewHandler(service auth.Service) *authHandler {
	return &authHandler{service: service}
}
