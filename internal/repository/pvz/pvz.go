package pvz
import (
	"avito-pvz/internal/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PvzRepository interface {
	CreatePVZ(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error)
	GetCityIdByName(ctx context.Context, city *entity.City) (int, error)
}

type pvzRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) PvzRepository {
	return &pvzRepository{pool: pool}
}
