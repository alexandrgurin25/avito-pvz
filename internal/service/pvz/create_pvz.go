package pvz

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"avito-pvz/pkg/logger"
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

func (s *pVZService) CreatePVZ(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error) {
	cityID, err := s.repo.GetCityIdByName(ctx, &pvz.City)
	if err != nil {
		if errors.Is(err, myerrors.ErrCityNotFound) {
			return nil, myerrors.ErrInvalidCity
		}
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to get city ID:", zap.Error(err))
		return nil, fmt.Errorf("failed to get city ID: %v", err)
	}

	pvz.City.Id = cityID

	pvz, err = s.repo.CreatePVZ(ctx, pvz)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to create PVZ",
			zap.Any("pvz", pvz),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create PVZ: %v", err)
	}

	return pvz, nil
}
