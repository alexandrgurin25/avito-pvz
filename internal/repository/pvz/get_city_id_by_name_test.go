package pvz

import (
	"avito-pvz/internal/config"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"avito-pvz/pkg/postgres"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCityIdByName(t *testing.T) {
	// Подготовка тестовых данных
	ctx := context.Background()
	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	cityName := "Москва"

	// Тест 1: Успешное получение ID города
	t.Run("Success", func(t *testing.T) {
		city := &entity.City{Name: cityName}
		cityID, err := repo.GetCityIdByName(ctx, city)
		assert.NoError(t, err)
		assert.NotZero(t, cityID) // Проверяем, что ID города не равен нулю
	})

	// Тест 2: Город не найден
	t.Run("CityNotFound", func(t *testing.T) {
		city := &entity.City{Name: "Неизвестный Город"}
		cityID, err := repo.GetCityIdByName(ctx, city)
		assert.Error(t, err)
		assert.Equal(t, myerrors.ErrCityNotFound, err)
		assert.Zero(t, cityID) // Проверяем, что ID города равен нулю
	})

	assert.NoError(t, err)
}
