package reception

import (
	"avito-pvz/internal/constants"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"context"
	"errors"
	"fmt"
)

func (s *receptionService) CreateReception(ctx context.Context, pvzId string) (*entity.Reception, error) {
	//Существует ли такое pvz

	_, err := s.pvzRepository.GetPvzById(ctx, pvzId)
	if err != nil {
		if errors.Is(err, myerrors.ErrPVZNotFound) {
			return nil, myerrors.ErrPVZNotFound
		}
		return nil, fmt.Errorf("failed to get pvz by ID: %v", err)
	}

	//Проверить, все ли приемки закрыты
	activeReception, err := s.receptionRepository.GetActiveReception(ctx, pvzId)
	if err != nil {
		return nil, fmt.Errorf("failed to get active reception: %v", err)
	}

	if activeReception != nil {
		return nil, myerrors.ErrActiveReceptionFound
	}

	var pvz = &entity.Reception{}

	pvz.PvzID = pvzId
	pvz.Status = constants.StatusReceptionInProgres

	//Создать приемку, если она не открыта
	pvz, err = s.receptionRepository.CreateReception(ctx, pvz)
	if err != nil {
		return nil, fmt.Errorf("failed to create reception: %v", err)
	}
	return pvz, nil
}
