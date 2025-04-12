package reception

import (
	"avito-pvz/internal/entity"
	"context"
	"database/sql"
	"errors"
)

func (r *receptionRepository) GetActiveReception(ctx context.Context, pvzID string) (*entity.Reception, error) {
	var reception entity.Reception

	err := r.db.QueryRow(
		ctx,
		`SELECT id, pvz_id, start_time, status 
		FROM receivings
		WHERE pvz_id = $1 AND status = 'in_progress'
		LIMIT 1`,
		pvzID,
	).Scan(&reception.ID, &reception.PvzID, &reception.DateTime, &reception.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &reception, nil
}
