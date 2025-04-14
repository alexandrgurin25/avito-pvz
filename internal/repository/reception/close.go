package reception

import (
	"avito-pvz/internal/entity"
	"context"
)

func (r *receptionRepository) CloseReception(ctx context.Context,
	reception *entity.Reception) (*entity.Reception, error) {

	var updated entity.Reception

	err := r.db.QueryRow(
		ctx,
		` UPDATE receivings
		SET status = $2, end_time = $3
		WHERE id = $1
		RETURNING id, pvz_id, start_time, status, end_time`,
		reception.ID,
		reception.Status,
		reception.CloseTime,
	).Scan(
		&updated.ID,
		&updated.PvzID,
		&updated.DateTime,
		&updated.Status,
		&updated.CloseTime,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}
