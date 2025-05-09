package pvz

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *pvzRepository) CreatePVZ(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error) {
	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO pvz (id, city_id, created_at) VALUES ($1, $2, $3)`,
		pvz.UUID,
		pvz.City.Id,
		pvz.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, myerrors.ErrPVZAlreadyExists
		}
		return nil, fmt.Errorf("falied pvz create error %v", err)
	}

	return &entity.PVZ{
		UUID: pvz.UUID,
		City: entity.City{
			Name: pvz.City.Name,
		},
		CreatedAt: pvz.CreatedAt,
	}, nil

}
