package pvz

import (
	"avito-pvz/internal/entity"
	"avito-pvz/internal/repository/pvz"
	"context"
)

type PVZService interface {
	CreatePVZ(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error)
}

type pVZService struct {
	repo pvz.Repository
}

func NewService(repo pvz.Repository) PVZService {
	return &pVZService{repo: repo}
}
