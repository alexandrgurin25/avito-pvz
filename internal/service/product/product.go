package product

import (
	"avito-pvz/internal/entity"
	"avito-pvz/internal/repository/product"
	"avito-pvz/internal/repository/pvz"
	"avito-pvz/internal/repository/reception"
	"context"
)

//go:generate mockgen -source=product.go -destination=mocks/product_mock.go -package=mocks
type Service interface {
	AddProduct(ctx context.Context, categoryId string, pvzId string) (*entity.Product, error)
	DeleteLastProduct(ctx context.Context, productId string) error
}

type productService struct {
	productRepository   product.Repository
	pvzRepository       pvz.Repository
	receptionRepository reception.Repository
}

func New(productRepository product.Repository,
	pvzRepository pvz.Repository,
	receptionRepository reception.Repository) Service {
	return &productService{productRepository: productRepository, pvzRepository: pvzRepository, receptionRepository: receptionRepository}
}
