package pvz

import (
	"avito-pvz/internal/entity"
	"avito-pvz/internal/repository/pvz"
	"context"
	"time"
)

type PVZService interface {
	CreatePVZ(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error)
	GetAllWithReceptions(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]entity.PVZ, error)
}

type pVZService struct {
	repo pvz.Repository
}

func NewService(repo pvz.Repository) PVZService {
	return &pVZService{repo: repo}
}
