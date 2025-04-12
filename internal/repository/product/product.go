package product

// import (
// 	"avito-pvz/internal/entity"
// 	"context"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type ProductRepository interface {
// 	AddProduct(ctx context.Context, product *entity.Product) error
// 	GetLastProduct(ctx context.Context, receptionID string) (*entity.Product, error)
// 	DeleteProduct(ctx context.Context, productID string) error
// 	GetProductsByReception(ctx context.Context, receptionID string) ([]entity.Product, error)
// }

// type productRepo struct {
// 	db *pgxpool.Pool
// }

// func NewRepository(db *pgxpool.Pool) ProductRepository {
// 	return &productRepo{db: db}
// }
