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

func TestCreateReception(t *testing.T) {
	// Подготовка тестовых данных
	ctx := context.Background()
	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	cityId := "1"
	pvzID := uuid.New().String()
	receivingID := uuid.New().String()

	// Создаем тестовый PVZ
	_, err = db.Exec(ctx, `INSERT INTO pvz (id, city_id) VALUES ($1, $2)`, pvzID, cityId)
	assert.NoError(t, err)

	// Создаем тестовый receiving
	_, err = db.Exec(ctx, `INSERT INTO receivings (id, pvz_id, status, start_time) VALUES ($1, $2, $3, $4)`,
		receivingID, pvzID, "in_progress", time.Now())
	assert.NoError(t, err)

	// Тест 1: Успешное создание получения
	t.Run("Success", func(t *testing.T) {
		reception := &entity.Reception{
			PvzID:  pvzID,
			Status: "in_progress",
		}

		createdReception, err := repo.CreateReception(ctx, reception)
		assert.NoError(t, err)
		assert.NotNil(t, createdReception)
		assert.Equal(t, reception.PvzID, createdReception.PvzID)
		assert.Equal(t, reception.Status, createdReception.Status)
		assert.NotEmpty(t, createdReception.ID)
		assert.NotZero(t, createdReception.DateTime)
	})

	// Тест 2: Ошибка при создании получения (например, если PVZ не существует)
	t.Run("ErrorCreatingReception", func(t *testing.T) {
		reception := &entity.Reception{
			PvzID:  uuid.New().String(), // Несуществующий ID PVZ
			Status: "in_progress",
		}

		createdReception, err := repo.CreateReception(ctx, reception)
		assert.Error(t, err)
		assert.Nil(t, createdReception)
	})

	// Очищаем тестовые данные

	_, err = db.Exec(ctx, `DELETE FROM receivings WHERE pvz_id = $1`, pvzID)
	assert.NoError(t, err)

}
