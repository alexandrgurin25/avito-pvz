package reception

import (
	"avito-pvz/internal/constants"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"avito-pvz/pkg/logger"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func (s *receptionService) CloseLastReception(ctx context.Context, pvzID string) (*entity.Reception, error) {
	// Находим активную приемку
	activeReception, err := s.receptionRepository.GetActiveReception(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active reception: %w", err)
	}
	if activeReception == nil {
		return nil, myerrors.ErrActiveReceptionNotFound
	}

	timeString := time.Now().Format(time.RFC3339)

	date, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "invalid parsing date", zap.Error(err))
	}

	// Закрываем приемку
	activeReception.Status = constants.StatusReceptionClose
	activeReception.CloseTime = date

	closedReception, err := s.receptionRepository.CloseReception(ctx, activeReception)
	if err != nil {
		return nil, fmt.Errorf("failed to close reception: %v", err)
	}
 
	closedReception.DateTime, err = time.Parse(time.RFC3339, closedReception.DateTime.String())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "invalid parsing date", zap.Error(err))
	}

	return closedReception, nil
}
