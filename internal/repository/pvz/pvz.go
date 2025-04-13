package pvz

import (
	"avito-pvz/internal/entity"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -source=pvz.go -destination=mocks/pvz_mock.go -package=mocks
type Repository interface {
	CreatePVZ(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error)
	GetCityIdByName(ctx context.Context, city *entity.City) (int, error)
	GetPvzById(ctx context.Context, uuid string) (*entity.PVZ, error)
	GetPVZsWithFilters(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entity.PVZ, error)
}

type pvzRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pvzRepository{pool: pool}
}
