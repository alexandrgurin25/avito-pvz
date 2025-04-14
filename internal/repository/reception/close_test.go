package reception

import (
	"avito-pvz/internal/config"
	"avito-pvz/internal/entity"
	"avito-pvz/pkg/postgres"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCloseReception(t *testing.T) {
	// Подготовка тестовых данных
	ctx := context.Background()
	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	// Generate unique test data
	cityId := 1
	pvzID := uuid.New().String()
	receivingID := uuid.New().String()

	// Создаем тестовый PVZ
	_, err = db.Exec(ctx, `INSERT INTO pvz (id, city_id) VALUES ($1, $2)`, pvzID, cityId)
	assert.NoError(t, err)

	// Создаем тестовый receiving
	_, err = db.Exec(ctx, `INSERT INTO receivings (id, pvz_id, status, start_time) VALUES ($1, $2, $3, $4)`,
		receivingID, pvzID, "in_progress", time.Now())
	assert.NoError(t, err)

	// Тест 1: Успешное закрытие получения
	t.Run("Success", func(t *testing.T) {
		reception := &entity.Reception{
			ID:     receivingID,
			PvzID:  pvzID,
			Status: "close",
		}

		closedReception, err := repo.CloseReception(ctx, reception)
		assert.NoError(t, err)
		assert.NotNil(t, closedReception)
		assert.Equal(t, reception.ID, closedReception.ID)
		assert.Equal(t, reception.Status, closedReception.Status)
	})

	// Тест 2: Ошибка при закрытии получения (например, если ID не существует)
	t.Run("ErrorClosingReception", func(t *testing.T) {
		reception := &entity.Reception{
			ID:        uuid.New().String(), // Несуществующий ID
			PvzID:     pvzID,
			Status:    "closed",
			CloseTime: time.Now(),
		}

		closedReception, err := repo.CloseReception(ctx, reception)
		assert.Error(t, err)
		assert.Nil(t, closedReception)
	})

	// Очищаем тестовые данные
	_, err = db.Exec(ctx, `DELETE FROM receivings WHERE id = $1`, receivingID)
	assert.NoError(t, err)
	_, err = db.Exec(ctx, `DELETE FROM pvz WHERE id = $1`, pvzID)
	assert.NoError(t, err)
}
