package pvz

import (
	"avito-pvz/internal/entity"
	"context"
	"fmt"
	"time"
)

func (s *pVZService) GetAllWithReceptions(
	ctx context.Context,
	startDate, endDate *time.Time,
	page, limit int,
) ([]entity.PVZ, error) {
	pvz, err := s.repo.GetPVZsWithFilters(ctx, startDate, endDate, page, limit)

	if err != nil {
		return nil, fmt.Errorf("failed to query pvz with receptions: %v", err)
	}

	return pvz, nil
}
