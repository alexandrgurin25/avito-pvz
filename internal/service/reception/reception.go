package reception

import (
	"avito-pvz/internal/entity"
	"avito-pvz/internal/repository/pvz"
	"avito-pvz/internal/repository/reception"
	"context"
)

//go:generate mockgen -source=reception.go -destination=mocks/reception.go -package=mocks

type Service interface {
	CreateReception(ctx context.Context, pvzID string) (*entity.Reception, error)
	CloseLastReception(ctx context.Context, pvzID string) (*entity.Reception, error)
}

type receptionService struct {
	receptionRepository reception.Repository
	pvzRepository       pvz.Repository
}

func NewService(receptionRepository reception.Repository, pvzRepository pvz.Repository) Service {
	return &receptionService{receptionRepository: receptionRepository, pvzRepository: pvzRepository}
}
