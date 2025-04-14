package pvz

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"avito-pvz/internal/repository/pvz/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreatePVZ_InvalidCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	inputPVZ := &entity.PVZ{
		UUID: "Test PVZ",
		City: entity.City{
			Name: "Non-existent City",
		},
	}

	mockRepo := mocks.NewMockRepository(ctrl)

	// Expectations
	mockRepo.EXPECT().
		GetCityIdByName(ctx, &inputPVZ.City).
		Return(0, myerrors.ErrCityNotFound)

	// Create service
	service := NewService(mockRepo)

	// Call method
	result, err := service.CreatePVZ(ctx, inputPVZ)

	// Verify
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, myerrors.ErrInvalidCity)
}

func TestCreatePVZ_GetCityIdError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	inputPVZ := &entity.PVZ{
		UUID: "Test PVZ",
		City: entity.City{
			Name: "Test City",
		},
	}
	dbErr := errors.New("database error")

	mockRepo := mocks.NewMockRepository(ctrl)

	// Expectations
	mockRepo.EXPECT().
		GetCityIdByName(ctx, &inputPVZ.City).
		Return(0, dbErr)

	// Create service
	service := NewService(mockRepo)

	// Call method
	result, err := service.CreatePVZ(ctx, inputPVZ)

	// Verify
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get city ID")
}

func TestCreatePVZ_PVZAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	inputPVZ := &entity.PVZ{
		UUID: "Existing PVZ",
		City: entity.City{
			Name: "Test City",
		},
	}
	expectedCityID := 789

	mockRepo := mocks.NewMockRepository(ctrl)

	// Expectations
	mockRepo.EXPECT().
		GetCityIdByName(ctx, &inputPVZ.City).
		Return(expectedCityID, nil)

	mockRepo.EXPECT().
		CreatePVZ(ctx, gomock.Any()).
		Return(nil, myerrors.ErrPVZAlreadyExists)

	// Create service
	service := NewService(mockRepo)

	// Call method
	result, err := service.CreatePVZ(ctx, inputPVZ)

	// Verify
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, myerrors.ErrPVZAlreadyExists)
}

func TestCreatePVZ_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	inputPVZ := &entity.PVZ{

		City: entity.City{
			Name: "Test City",
		},
	}
	expectedCityID := 4
	expectedPVZ := &entity.PVZ{
		UUID: "3fa85f64-5717-4562-b3fc-2c963f66afa6",

		City: entity.City{
			Id:   expectedCityID,
			Name: inputPVZ.City.Name,
		},
	}

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		GetCityIdByName(ctx, &inputPVZ.City).
		Return(expectedCityID, nil)

	mockRepo.EXPECT().
		CreatePVZ(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, pvz *entity.PVZ) (*entity.PVZ, error) {
			assert.Equal(t, expectedCityID, pvz.City.Id)
			return expectedPVZ, nil
		})

	servic := NewService(mockRepo)

	// Call method
	result, err := servic.CreatePVZ(ctx, inputPVZ)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, expectedPVZ, result)
}

func TestCreatePVZ_CreatePVZError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	inputPVZ := &entity.PVZ{
		UUID: "Test PVZ",
		City: entity.City{
			Name: "Test City",
		},
	}
	expectedCityID := 101
	dbErr := errors.New("database error")

	mockRepo := mocks.NewMockRepository(ctrl)

	// Expectations
	mockRepo.EXPECT().
		GetCityIdByName(ctx, &inputPVZ.City).
		Return(expectedCityID, nil)

	mockRepo.EXPECT().
		CreatePVZ(ctx, gomock.Any()).
		Return(nil, dbErr)

	// Create service
	service := NewService(mockRepo)

	// Call method
	result, err := service.CreatePVZ(ctx, inputPVZ)

	// Verify
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create PVZ")
}
