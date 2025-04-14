package reception

import (
	"avito-pvz/internal/entity"
	"context"
	"fmt"
)

func (r *receptionRepository) CreateReception(
	ctx context.Context,
	reception *entity.Reception) (*entity.Reception, error) {

	err := r.db.QueryRow(
		ctx,
		`INSERT INTO receivings (pvz_id, status)
		VALUES ($1, $2)
		RETURNING id, start_time`,
		reception.PvzID,
		reception.Status,
	).Scan(&reception.ID, &reception.DateTime)

	if err != nil {
		return nil, fmt.Errorf("failed to create reception: %v", err)
	}

	return &entity.Reception{
		ID:       reception.ID,
		PvzID:    reception.PvzID,
		DateTime: reception.DateTime,
		Status:   reception.Status,
	}, nil
}
