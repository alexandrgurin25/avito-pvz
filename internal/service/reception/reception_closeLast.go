package reception

import (
	"avito-pvz/internal/constants"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"context"
	"fmt"
	"time"
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

	// Закрываем приемку
	activeReception.Status = constants.StatusReceptionClose
	activeReception.CloseTime = time.Now().Format("2025-04-12T19:20:48.201Z")

	closedReception, err := s.receptionRepository.UpdateReception(ctx, activeReception)
	if err != nil {
		return nil, fmt.Errorf("failed to close reception: %v", err)
	}

	return closedReception, nil
}
