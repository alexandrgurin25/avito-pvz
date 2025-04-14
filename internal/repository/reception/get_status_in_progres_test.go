package reception

import (
	"avito-pvz/internal/config"
	"avito-pvz/pkg/postgres"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetActiveReception(t *testing.T) {
	// Подготовка тестовых данных
	ctx := context.Background()
	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	pvzID := uuid.New().String()
	receivingID := uuid.New().String()

	// Создаем тестовый город
	var cityId = 1

	// Создаем тестовый PVZ
	_, err = db.Exec(ctx, `INSERT INTO pvz (id, city_id) VALUES ($1, $2)`, pvzID, cityId)
	assert.NoError(t, err)

	// Создаем тестовый receiving
	_, err = db.Exec(ctx, `INSERT INTO receivings (id, pvz_id, status, start_time) VALUES ($1, $2, $3, $4)`,
		receivingID, pvzID, "in_progress", time.Now())
	assert.NoError(t, err)

	// Тест 1: Успешное получение активного получения
	t.Run("Success", func(t *testing.T) {
		activeReception, err := repo.GetActiveReception(ctx, pvzID)
		assert.NoError(t, err)
		assert.NotNil(t, activeReception)
		assert.Equal(t, receivingID, activeReception.ID)
		assert.Equal(t, pvzID, activeReception.PvzID)
		assert.Equal(t, "in_progress", activeReception.Status)
	})

	// Тест 2: Ошибка при получении активного получения (например, если активного получения не существует)
	t.Run("NoActiveReception", func(t *testing.T) {
		// Удаляем активное получение
		_, err := db.Exec(ctx, `UPDATE receivings SET status = 'close' WHERE id = $1`, receivingID)
		assert.NoError(t, err)

		activeReception, err := repo.GetActiveReception(ctx, pvzID)
		assert.NoError(t, err)
		assert.Nil(t, activeReception)
	})

	// Очищаем тестовые данные
	_, err = db.Exec(ctx, `DELETE FROM receivings WHERE id = $1`, receivingID)
	assert.NoError(t, err)
	_, err = db.Exec(ctx, `DELETE FROM pvz WHERE id = $1`, pvzID)
	assert.NoError(t, err)
}
