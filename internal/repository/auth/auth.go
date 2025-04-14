package auth

import (
	"avito-pvz/internal/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, email string, passwordHash string, role string) (*entity.User, error)
}

type authRepo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &authRepo{db: db}
}
