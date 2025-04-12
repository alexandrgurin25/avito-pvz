package pvz

import (
	"avito-pvz/internal/entity"
	pvzR "avito-pvz/internal/repository/pvz"
	"context"
)

type PVZService interface {
	CreatePVZ(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error)
}

type pVZService struct {
	repo pvzR.PvzRepository
}

func NewService(repo pvzR.PvzRepository) PVZService {
	return &pVZService{repo: repo}
}
