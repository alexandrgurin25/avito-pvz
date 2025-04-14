package reception

import (
	"avito-pvz/internal/constants"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	pvz_mocks "avito-pvz/internal/repository/pvz/mocks"
	"avito-pvz/internal/repository/reception/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCloseLastReception_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	now := time.Now()

	mockRepo := mocks.NewMockRepository(ctrl)

	// Ожидания
	activeReception := &entity.Reception{
		ID:     "4fa85f64-5717-4562-b3fc-2c963f66afa6",
		PvzID:  pvzID,
		Status: constants.StatusReceptionInProgres,
	}
	mockRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(activeReception, nil)

	expectedClosedReception := &entity.Reception{
		ID:        activeReception.ID,
		PvzID:     pvzID,
		Status:    constants.StatusReceptionClose,
		CloseTime: now,
	}
	mockRepo.EXPECT().
		CloseReception(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, r *entity.Reception) (*entity.Reception, error) {
			assert.Equal(t, constants.StatusReceptionClose, r.Status)
			return expectedClosedReception, nil
		})

	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	// Создаем сервис
	service := NewService(mockRepo, mockPvzRepo)

	// Вызываем метод
	result, err := service.CloseLastReception(ctx, pvzID)

	// Проверяем
	require.NoError(t, err)
	assert.Equal(t, expectedClosedReception, result)
	assert.Equal(t, constants.StatusReceptionClose, result.Status)
	assert.WithinDuration(t, now, result.CloseTime, time.Second)
}

func TestCloseLastReception_GetActiveReceptionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	dbErr := errors.New("database error")

	mockRepo := mocks.NewMockRepository(ctrl)

	// Ожидания
	mockRepo.EXPECT().GetActiveReception(ctx, pvzID).
		Return(nil, dbErr)

	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	// Создаем сервис
	service := NewService(mockRepo, mockPvzRepo)

	// Вызываем метод
	result, err := service.CloseLastReception(ctx, pvzID)

	// Проверяем
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get active reception")
}

func TestCloseLastReception_ActiveReceptionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"

	mockRepo := mocks.NewMockRepository(ctrl)

	// Ожидания
	mockRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(nil, nil)

	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	// Создаем сервис
	service := NewService(mockRepo, mockPvzRepo)

	// Вызываем метод
	result, err := service.CloseLastReception(ctx, pvzID)

	// Проверяем
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, myerrors.ErrActiveReceptionNotFound)
}

func TestCloseLastReception_CloseReceptionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	dbErr := errors.New("database error")

	mockRepo := mocks.NewMockRepository(ctrl)

	// Ожидания
	activeReception := &entity.Reception{
		ID:     "4fa85f64-5717-4562-b3fc-2c963f66afa6",
		PvzID:  pvzID,
		Status: constants.StatusReceptionInProgres,
	}
	mockRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(activeReception, nil)

	mockRepo.EXPECT().
		CloseReception(ctx, gomock.Any()).
		Return(nil, dbErr)

	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	// Создаем сервис
	service := NewService(mockRepo, mockPvzRepo)

	// Вызываем метод
	result, err := service.CloseLastReception(ctx, pvzID)

	// Проверяем
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to close reception")
}
