package reception

import (
	"avito-pvz/internal/constants"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"avito-pvz/internal/repository/pvz/mocks"
	reception_mocks "avito-pvz/internal/repository/reception/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateReception_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"

	mockPvzRepo := mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания для PVZ
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	// Ожидания для Reception
	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(nil, nil)

	expectedReception := &entity.Reception{
		PvzID:  pvzID,
		Status: constants.StatusReceptionInProgres,
	}
	mockReceptionRepo.EXPECT().
		CreateReception(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, r *entity.Reception) (*entity.Reception, error) {
			assert.Equal(t, pvzID, r.PvzID)
			assert.Equal(t, constants.StatusReceptionInProgres, r.Status)
			return expectedReception, nil
		})

	// Создаем сервис
	service := NewService(mockReceptionRepo, mockPvzRepo)

	// Вызываем метод
	result, err := service.CreateReception(ctx, pvzID)

	// Проверяем
	require.NoError(t, err)
	assert.Equal(t, expectedReception, result)
}

func TestCreateReception_PVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"

	mockPvzRepo := mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания для PVZ
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(nil, myerrors.ErrPVZNotFound)

	// Создаем сервис
	service := NewService(mockReceptionRepo, mockPvzRepo)

	// Вызываем метод
	result, err := service.CreateReception(ctx, pvzID)

	// Проверяем
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, myerrors.ErrPVZNotFound)
}

func TestCreateReception_GetPVZError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	dbErr := errors.New("database error")

	mockPvzRepo := mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания для PVZ
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(nil, dbErr)

	// Создаем сервис
	service := NewService(mockReceptionRepo, mockPvzRepo)

	// Вызываем метод
	result, err := service.CreateReception(ctx, pvzID)

	// Проверяем
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get pvz by ID")
	assert.Contains(t, err.Error(), dbErr.Error())
}

func TestCreateReception_ActiveReceptionFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"

	mockPvzRepo := mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания для PVZ
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	// Ожидания для Reception
	activeReception := &entity.Reception{PvzID: pvzID, Status: constants.StatusReceptionInProgres}
	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(activeReception, nil)

	// Создаем сервис
	service := NewService(mockReceptionRepo, mockPvzRepo)

	// Вызываем метод
	result, err := service.CreateReception(ctx, pvzID)

	// Проверяем
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, myerrors.ErrActiveReceptionFound)
}

func TestCreateReception_GetActiveReceptionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	dbErr := errors.New("database error")

	mockPvzRepo := mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания для PVZ
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	// Ожидания для Reception
	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(nil, dbErr)

	// Создаем сервис
	service := NewService(mockReceptionRepo, mockPvzRepo)

	// Вызываем метод
	result, err := service.CreateReception(ctx, pvzID)

	// Проверяем
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get active reception")
	assert.Contains(t, err.Error(), dbErr.Error())
}

func TestCreateReception_CreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	dbErr := errors.New("database error")

	mockPvzRepo := mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания для PVZ
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	// Ожидания для Reception
	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(nil, nil)

	mockReceptionRepo.EXPECT().
		CreateReception(ctx, gomock.Any()).
		Return(nil, dbErr)

	// Создаем сервис
	service := NewService(mockReceptionRepo, mockPvzRepo)

	// Вызываем метод
	result, err := service.CreateReception(ctx, pvzID)

	// Проверяем
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create reception")
	assert.Contains(t, err.Error(), dbErr.Error())
}
