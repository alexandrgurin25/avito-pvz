package pvz

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"context"
	"database/sql"
	"errors"
)

func (r *pvzRepository) GetPvzById(ctx context.Context, uuid string) (*entity.PVZ, error) {
	var pvz entity.PVZ

	err := r.pool.QueryRow(
		ctx,
		`SELECT id, city_id, created_at FROM pvz WHERE id = $1`,
		uuid,
	).Scan(&pvz.UUID, &pvz.City.Id, &pvz.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.ErrPVZNotFound
		}
		return nil, err
	}
	return &pvz, nil
}
