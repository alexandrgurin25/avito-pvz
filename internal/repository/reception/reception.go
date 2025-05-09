package reception

import (
	"avito-pvz/internal/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -source=reception.go -destination=mocks/reception_mock.go -package=mocks
type Repository interface {
	CreateReception(ctx context.Context, reception *entity.Reception) (*entity.Reception, error)
	GetActiveReception(ctx context.Context, pvzID string) (*entity.Reception, error)
	CloseReception(ctx context.Context, reception *entity.Reception) (*entity.Reception, error)
}

type receptionRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &receptionRepository{db: db}
}
