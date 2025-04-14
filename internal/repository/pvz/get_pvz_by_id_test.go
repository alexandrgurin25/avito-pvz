package pvz

import (
	"avito-pvz/internal/config"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/pkg/postgres"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetPvzById(t *testing.T) {
	// Подготовка тестовых данных
	ctx := context.Background()
	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	// Создаем тестовый PVZ
	pvzID := uuid.New().String()
	cityID := 1 
	_, err = db.Exec(ctx, `INSERT INTO pvz (id, city_id, created_at) VALUES ($1, $2, $3)`, pvzID, cityID, "2023-01-01 00:00:00")
	assert.NoError(t, err)

	// Тест 1: Успешное получение PVZ по ID
	t.Run("Success", func(t *testing.T) {
		pvz, err := repo.GetPvzById(ctx, pvzID)
		assert.NoError(t, err)
		assert.NotNil(t, pvz)                // Проверяем, что PVZ не равен nil
		assert.Equal(t, pvzID, pvz.UUID)     // Проверяем, что UUID совпадает
		assert.Equal(t, cityID, pvz.City.Id) // Проверяем, что city_id совпадает
	})

	// Тест 2: PVZ не найден
	t.Run("PVZNotFound", func(t *testing.T) {
		invalidID := uuid.New().String() // Генерируем новый UUID, который не существует в базе
		pvz, err := repo.GetPvzById(ctx, invalidID)
		assert.Error(t, err)
		assert.Equal(t, myerrors.ErrPVZNotFound, err) // Проверяем, что ошибка соответствует ожидаемой
		assert.Nil(t, pvz)                            // Проверяем, что PVZ равен nil
	})

	// Очищаем тестовые данные
	_, err = db.Exec(ctx, `DELETE FROM pvz WHERE id = $1`, pvzID)
	assert.NoError(t, err)
}
