package test

import (
	"avito-pvz/internal/config"
	"avito-pvz/internal/constants"
	"avito-pvz/internal/entity"
	"avito-pvz/internal/repository/product"
	"avito-pvz/internal/repository/pvz"
	"avito-pvz/internal/repository/reception"
	productService "avito-pvz/internal/service/product"
	pvzService "avito-pvz/internal/service/pvz"
	receptionService "avito-pvz/internal/service/reception"
	"avito-pvz/pkg/postgres"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFullPVZWorkflow(t *testing.T) {
	ctx := context.Background()

	cfg, err := config.NewTest("../../config/test.env")
	require.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, "file://../../db/migrations")
	require.NoError(t, err, "unable to connect db")

	// 1. Создаем новый ПВЗ
	pvzRepo := pvz.NewRepository(db)
	pvzSvc := pvzService.NewService(pvzRepo)

	city := "Москва"
	newPVZ := &entity.PVZ{
		UUID: uuid.New().String(),
		City: entity.City{
			Name: city,
		},
		CreatedAt: time.Now(),
	}

	createdPVZ, err := pvzSvc.CreatePVZ(ctx, newPVZ)
	require.NoError(t, err)
	require.NotEmpty(t, createdPVZ.UUID)

	// 2. Создаем новую приемку
	receptionRepo := reception.NewRepository(db) // Используем db вместо s.pool
	receptionSvc := receptionService.NewService(receptionRepo, pvzRepo)

	createdReception, err := receptionSvc.CreateReception(ctx, createdPVZ.UUID)
	require.NoError(t, err)
	require.Equal(t, constants.StatusReceptionInProgres, createdReception.Status)

	// 3. Добавляем 50 товаров
	productRepo := product.NewRepository(db) // Используем db вместо s.pool
	productSvc := productService.New(productRepo, pvzRepo, receptionRepo)

	productTypes := []string{"электроника", "одежда", "обувь"}
	for i := 0; i < 50; i++ {
		productType := productTypes[i%3] // чередуем типы товаров
		_, err := productSvc.AddProduct(ctx, productType, createdPVZ.UUID)
		require.NoError(t, err)
	}

	// 4. Закрываем приемку
	closedReception, err := receptionSvc.CloseLastReception(ctx, createdPVZ.UUID)
	require.NoError(t, err)
	require.Equal(t, constants.StatusReceptionClose, closedReception.Status)
	require.False(t, closedReception.CloseTime.IsZero())
}
