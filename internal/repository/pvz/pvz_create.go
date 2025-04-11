package pvz

import (
	"avito-pvz/internal/entity"
	"context"
	"fmt"
	"time"
)

func (r *pvzRepository) CreatePVZ(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error) {
	var timeNow time.Time

	err := r.pool.QueryRow(
		ctx,
		`INSERT INTO pvz(id, city) VALUES $1, $2 RETURNING "created_at"`,
		pvz.UUID,
		pvz.City.Id,
	).Scan(&timeNow)

	if err != nil {
		return nil, fmt.Errorf("falied pvz create error %v", err)
	}

	return &entity.PVZ{
		CreatedAt: timeNow,
	}, nil

}
