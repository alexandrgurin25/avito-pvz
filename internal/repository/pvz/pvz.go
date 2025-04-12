package pvz

import (
	"avito-pvz/internal/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreatePVZ(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error)
	GetCityIdByName(ctx context.Context, city *entity.City) (int, error)
	GetPvzById(ctx context.Context, uuid string) (*entity.PVZ, error)
}

type pvzRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pvzRepository{pool: pool}
}
