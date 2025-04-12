package reception

import (
	"avito-pvz/internal/entity"
	"avito-pvz/internal/repository/pvz"
	"avito-pvz/internal/repository/reception"
	"context"
)

type Service interface {
	CreateReception(ctx context.Context, pvzID string) (*entity.Reception, error)
	// CloseReception(ctx context.Context, pvzID string) (*entity.Reception, error)
}

type receptionService struct {
	receptionRepository reception.Repository
	pvzRepository       pvz.Repository
}

func NewService(receptionRepository reception.Repository, pvzRepository pvz.Repository) Service {
	return &receptionService{receptionRepository: receptionRepository, pvzRepository: pvzRepository}
}
