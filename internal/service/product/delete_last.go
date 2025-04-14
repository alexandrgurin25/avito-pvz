package product

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/pkg/logger"
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

func (s *productService) DeleteLastProduct(ctx context.Context, pvzID string) error {
	// Проверим существует ли такой pvz
	_, err := s.pvzRepository.GetPvzById(ctx, pvzID)
	if err != nil {
		if errors.Is(err, myerrors.ErrPVZNotFound) {
			return myerrors.ErrPVZNotFound
		}
		return fmt.Errorf("failed to get pvz by ID: %v", err)
	}

	// Получаем активную приемку
	activeReception, err := s.receptionRepository.GetActiveReception(ctx, pvzID)
	if err != nil {
		return fmt.Errorf("failed to get active reception: %v", err)
	}
	if activeReception == nil {
		return myerrors.ErrActiveReceptionNotFound
	}

	//Находим id последнего товара
	product, err := s.productRepository.GetLastProductByReceigingId(ctx, activeReception.ID)
	if err != nil {
		if errors.Is(err, myerrors.ErrNoProductsToDelete) {
			return err
		}
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get last product:", zap.Error(err))
		return fmt.Errorf("failed to get last product: %v", err)
	}

	// Удаляем последний товар
	err = s.productRepository.DeleteLastProduct(ctx, product.ID)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to delete last product:", zap.Error(err))
		return fmt.Errorf("failed to delete last product: %v", err)
	}

	return nil
}
