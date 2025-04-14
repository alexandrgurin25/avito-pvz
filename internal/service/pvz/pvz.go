package pvz

import (
	"avito-pvz/internal/entity"
	"avito-pvz/internal/repository/pvz"
	"context"
	"time"
)

//go:generate mockgen -source=pvz.go -destination=mocks/pvz_mock.go -package=mocks
type PVZService interface {
	CreatePVZ(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error)
	GetAllWithReceptions(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entity.PVZ, error)
}

type pVZService struct {
	repo      pvz.Repository
	cityCache *cache
}

func NewService(repo pvz.Repository) PVZService {
	cache := NewCache()
	return &pVZService{repo: repo, cityCache: cache}
}
