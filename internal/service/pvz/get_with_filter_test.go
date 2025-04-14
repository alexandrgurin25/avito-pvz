package pvz

import (
	"avito-pvz/internal/entity"
	"avito-pvz/internal/repository/pvz/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAllWithReceptions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
	page := 1
	limit := 10

	expectedPVZs := []entity.PVZ{
		{
			UUID: "1fa85f64-5717-4562-b3fc-2c963f66afa6",
			City: entity.City{
				Id:   1,
				Name: "City 1",
			},
		},
		{
			UUID: "2fa85f64-5717-4562-b3fc-2c963f66afa6",
			City: entity.City{
				Id:   2,
				Name: "City 2",
			},
		},
	}

	mockRepo := mocks.NewMockRepository(ctrl)

	// Expectations
	mockRepo.EXPECT().
		GetPVZsWithFilters(ctx, &startDate, &endDate, page, limit).
		Return(expectedPVZs, nil)

	// Create service
	service := NewService(mockRepo)

	// Call method
	result, err := service.GetAllWithReceptions(ctx, &startDate, &endDate, page, limit)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, expectedPVZs, result)
}

func TestGetAllWithReceptions_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	startDate := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 6, 30, 23, 59, 59, 0, time.UTC)
	page := 2
	limit := 20
	dbErr := errors.New("database error")

	mockRepo := mocks.NewMockRepository(ctrl)

	// Expectations
	mockRepo.EXPECT().
		GetPVZsWithFilters(ctx, &startDate, &endDate, page, limit).
		Return(nil, dbErr)

	// Create service
	service := NewService(mockRepo)

	// Call method
	result, err := service.GetAllWithReceptions(ctx, &startDate, &endDate, page, limit)

	// Verify
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to query pvz with receptions")
}
