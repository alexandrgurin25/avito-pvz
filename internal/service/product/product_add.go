package product

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"avito-pvz/pkg/logger"
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

func (s *productService) AddProduct(
	ctx context.Context,
	categoryName string,
	pvzId string) (*entity.Product, error) {

	//Проверим существует ли такой pvz
	_, err := s.pvzRepository.GetPvzById(ctx, pvzId)
	if err != nil {
		if errors.Is(err, myerrors.ErrPVZNotFound) {
			return nil, myerrors.ErrPVZNotFound
		}
		return nil, fmt.Errorf("failed to get pvz by ID: %v", err)
	}

	//Существует ли активная приемка
	activeReception, err := s.receptionRepository.GetActiveReception(ctx, pvzId)
	if err != nil {
		return nil, fmt.Errorf("failed to get active reception: %v", err)
	}
	if activeReception == nil {
		return nil, myerrors.ErrActiveReceptionNotFound
	}

	//Получить type_id по name type
	categoryId, err := s.productRepository.GetIdCategoryByName(ctx, categoryName)
	if err != nil {
		if errors.Is(err, myerrors.ErrInvalidProductType) {
			return nil, err
		}
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get id category by name:", zap.Error(err))
		return nil, fmt.Errorf("failed to get id category by name: %v", err)
	}
	//Добавить продукт
	product, err := s.productRepository.AddProduct(ctx, pvzId, categoryId)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to add product in reception:", zap.Error(err))
		return nil, fmt.Errorf("failed to add product in reception: %v", err)
	}

	product.Category = categoryName

	return product, nil

}
