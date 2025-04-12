package product

import (
	"avito-pvz/internal/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	AddProduct(ctx context.Context, receiving_id string, categoryId int) (*entity.Product, error)
	GetIdCategoryByName(ctx context.Context, name string) (int, error)
	// GetLastProduct(ctx context.Context, receptionID string) (*entity.Product, error)
	// DeleteProduct(ctx context.Context, productID string) error
	// GetProductsByReception(ctx context.Context, receptionID string) ([]entity.Product, error)
}

type productRepo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &productRepo{db: db}
}
