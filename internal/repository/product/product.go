package product

import (
	"avito-pvz/internal/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -source=product.go -destination=mocks/product_mock.go -package=mocks

type Repository interface {
	AddProduct(ctx context.Context, receiving_id string, categoryId int) (*entity.Product, error)
	GetIdCategoryByName(ctx context.Context, name string) (int, error)
	DeleteLastProduct(ctx context.Context, productId string) error
	GetLastProductByReceigingId(ctx context.Context, receptionID string) (*entity.Product, error)
}

type productRepo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &productRepo{db: db}
}
