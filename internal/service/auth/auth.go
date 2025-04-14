package auth

import (
	"avito-pvz/internal/repository/auth"
	"context"
)

type Service interface {
	CreateDummyLogin(role string) (string, error)
	Login(ctx context.Context, username string, password string) (string, error)
	Register(ctx context.Context, email string, password string, role string) error
}

type authService struct {
	authRepository auth.Repository
}

func NewService(authRepository auth.Repository) Service {
	return &authService{authRepository: authRepository}
}
