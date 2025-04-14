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

func TestAddProduct_GetIdCategoryByName_ReturnsErrInvalidProductType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	category := "несуществующая_категория"
	receptionID := "4fa85f64-5717-4562-b3fc-2c963f66afa6"

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания для PVZ и Reception (успешные)
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(&entity.Reception{ID: receptionID, PvzID: pvzID}, nil)

	// Ожидание для категории - возвращаем ошибку
	mockProductRepo.EXPECT().
		GetIdCategoryByName(ctx, category).
		Return(0, myerrors.ErrInvalidProductType)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	result, err := service.AddProduct(ctx, category, pvzID)

	// Проверяем
	require.ErrorIs(t, err, myerrors.ErrInvalidProductType)
	require.Nil(t, result)

	// Проверяем, что AddProduct не вызывался
	mockProductRepo.EXPECT().AddProduct(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
}

func TestAddProduct_GetIdCategoryByName_ReturnsUnknownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	category := "электроника"
	receptionID := "4fa85f64-5717-4562-b3fc-2c963f66afa6"
	unknownErr := errors.New("some database error")
	actualErr := errors.New("failed to get id category by name: some database error")

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания для PVZ и Reception (успешные)
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(&entity.Reception{ID: receptionID, PvzID: pvzID}, nil)

	// Ожидание для категории - возвращаем произвольную ошибку
	mockProductRepo.EXPECT().
		GetIdCategoryByName(ctx, category).
		Return(0, unknownErr)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	result, err := service.AddProduct(ctx, category, pvzID)

	// Проверяем
	require.Error(t, err)
	require.Equal(t, actualErr, err)
	require.Nil(t, result)

	// Проверяем, что AddProduct не вызывался
	mockProductRepo.EXPECT().AddProduct(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
}

func TestAddProduct_GetPvzById_ReturnsUnknownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	category := "электроника"
	unknownErr := errors.New("some database error")
	actualErr := errors.New("failed to get pvz by ID: some database error")

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания - GetPvzById возвращает произвольную ошибку
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(nil, unknownErr)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	result, err := service.AddProduct(ctx, category, pvzID)

	// Проверяем
	assert.Error(t, err)
	assert.Equal(t, actualErr, err)
	assert.Nil(t, result)

	// Проверяем, что другие методы не вызывались
	ctrl.Finish()
}

func TestAddProduct_GetPvzById_ReturnsErrPVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	category := "электроника"

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания - GetPvzById возвращает ошибку
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(nil, myerrors.ErrPVZNotFound)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	result, err := service.AddProduct(ctx, category, pvzID)

	// Проверяем
	require.ErrorIs(t, err, myerrors.ErrPVZNotFound)
	require.Nil(t, result)

	// Проверяем, что другие методы не вызывались
	ctrl.Finish()
}

func TestAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	category := "электроника"
	categoryID := 1
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
		GetIdCategoryByName(ctx, category).
		Return(categoryID, nil)

	expectedProduct := &entity.Product{
		ID:          "1",
		ReceptionID: receptionID,
		Category:    category,
	}
	mockProductRepo.EXPECT().
		AddProduct(ctx, receptionID, categoryID).
		Return(expectedProduct, nil)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	result, err := service.AddProduct(ctx, category, pvzID)

	// Проверяем
	require.NoError(t, err)
	require.Equal(t, expectedProduct, result)
}

func TestAddProduct_GetActiveReception_ReturnsErrActiveReceptionNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	category := "электроника"

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
	result, err := service.AddProduct(ctx, category, pvzID)

	// Проверяем
	require.ErrorIs(t, err, myerrors.ErrActiveReceptionNotFound)
	require.Nil(t, result)

	// Проверяем, что следующие методы не вызывались
	mockProductRepo.EXPECT().GetIdCategoryByName(gomock.Any(), gomock.Any()).Times(0)
	mockProductRepo.EXPECT().AddProduct(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
}

func TestAddProduct_GetActiveReception_ReturnsUnknownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	category := "электроника"
	unknownErr := errors.New("database connection failed")
	actualErr := errors.New("failed to get active reception: database connection failed")

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(nil, unknownErr)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	result, err := service.AddProduct(ctx, category, pvzID)

	// Проверяем
	require.Error(t, err)
	require.Equal(t, actualErr, err)
	require.Nil(t, result)

	// Проверяем, что следующие методы не вызывались
	mockProductRepo.EXPECT().GetIdCategoryByName(gomock.Any(), gomock.Any()).Times(0)
	mockProductRepo.EXPECT().AddProduct(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
}

func TestAddProduct_AddProduct_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	category := "электроника"
	categoryID := 1
	receptionID := "4fa85f64-5717-4562-b3fc-2c963f66afa6"
	addProductErr := errors.New("failed to add product")
	actualErr := errors.New("failed to add product in reception: failed to add product")

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockPvzRepo := pvz_mocks.NewMockRepository(ctrl)
	mockReceptionRepo := reception_mocks.NewMockRepository(ctrl)

	// Ожидания для зависимых вызовов (успешные)
	mockPvzRepo.EXPECT().
		GetPvzById(ctx, pvzID).
		Return(&entity.PVZ{UUID: pvzID}, nil)

	mockReceptionRepo.EXPECT().
		GetActiveReception(ctx, pvzID).
		Return(&entity.Reception{ID: receptionID, PvzID: pvzID}, nil)

	mockProductRepo.EXPECT().
		GetIdCategoryByName(ctx, category).
		Return(categoryID, nil)

	// Ожидание для AddProduct - возвращаем ошибку
	mockProductRepo.EXPECT().
		AddProduct(ctx, receptionID, categoryID).
		Return(nil, addProductErr)

	// Создаем сервис
	service := New(mockProductRepo, mockPvzRepo, mockReceptionRepo)

	// Вызываем метод
	result, err := service.AddProduct(ctx, category, pvzID)

	// Проверяем
	require.Error(t, err)
	require.Equal(t, actualErr, err)
	require.Nil(t, result)
}

