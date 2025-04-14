package product

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"avito-pvz/internal/repository/product/mocks"
	pvz_mocks "avito-pvz/internal/repository/pvz/mocks"
	reception_mocks "avito-pvz/internal/repository/reception/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteLastProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	receptionID := "4fa85f64-5717-4562-b3fc-2c963f66afa6"
	productID := "5fa85f64-5717-4562-b3fc-2c963f66afa6"

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(&entity.Reception{ID: receptionID, PvzID: pvzID}, nil)

	expectedProduct := &entity.Product{
		ID:          productID,
		ReceptionID: receptionID,
	}
	mockProductRepo.EXPECT().
		GetLastProductByReceigingId(ctx, receptionID).
		Return(expectedProduct, nil)

	mockProductRepo.EXPECT().
		DeleteLastProduct(ctx, productID).
		Return(nil)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	err := service.DeleteLastProduct(ctx, pvzID)

	// Проверяем
	require.NoError(t, err)
}

func TestDeleteLastProduct_PVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(nil, myerrors.ErrPVZNotFound)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	err := service.DeleteLastProduct(ctx, pvzID)

	// Проверяем
	assert.ErrorIs(t, err, myerrors.ErrPVZNotFound)
}

func TestDeleteLastProduct_GetPVZError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	dbErr := errors.New("database error")

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(nil, dbErr)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	err := service.DeleteLastProduct(ctx, pvzID)

	// Проверяем
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get pvz by ID")
}

func TestDeleteLastProduct_ActiveReceptionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(nil, nil)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	err := service.DeleteLastProduct(ctx, pvzID)

	// Проверяем
	assert.ErrorIs(t, err, myerrors.ErrActiveReceptionNotFound)
}

func TestDeleteLastProduct_GetActiveReceptionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	dbErr := errors.New("database error")

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(nil, dbErr)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	err := service.DeleteLastProduct(ctx, pvzID)

	// Проверяем
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get active reception")
}

func TestDeleteLastProduct_NoProductsToDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	receptionID := "4fa85f64-5717-4562-b3fc-2c963f66afa6"

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(&entity.Reception{ID: receptionID, PvzID: pvzID}, nil)

	mockProductRepo.EXPECT().
		GetLastProductByReceigingId(ctx, receptionID).
		Return(nil, myerrors.ErrNoProductsToDelete)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	err := service.DeleteLastProduct(ctx, pvzID)

	// Проверяем
	assert.ErrorIs(t, err, myerrors.ErrNoProductsToDelete)
}

func TestDeleteLastProduct_GetLastProductError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	receptionID := "4fa85f64-5717-4562-b3fc-2c963f66afa6"
	dbErr := errors.New("database error")

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(&entity.Reception{ID: receptionID, PvzID: pvzID}, nil)

	mockProductRepo.EXPECT().
		GetLastProductByReceigingId(ctx, receptionID).
		Return(nil, dbErr)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	err := service.DeleteLastProduct(ctx, pvzID)

	// Проверяем
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get last product")
}

func TestDeleteLastProduct_DeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	receptionID := "4fa85f64-5717-4562-b3fc-2c963f66afa6"
	productID := "5fa85f64-5717-4562-b3fc-2c963f66afa6"
	dbErr := errors.New("database error")

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(&entity.Reception{ID: receptionID, PvzID: pvzID}, nil)

	expectedProduct := &entity.Product{
		ID:          productID,
		ReceptionID: receptionID,
	}
	mockProductRepo.EXPECT().
		GetLastProductByReceigingId(ctx, receptionID).
		Return(expectedProduct, nil)

	mockProductRepo.EXPECT().
		DeleteLastProduct(ctx, productID).
		Return(dbErr)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	err := service.DeleteLastProduct(ctx, pvzID)

	// Проверяем
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete last product")
}