package pvz

import (
	"avito-pvz/internal/config"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"avito-pvz/pkg/postgres"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatePVZ(t *testing.T) {
	// Подготовка тестовых данных
	ctx := context.Background()
	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	testUUID := uuid.New().String()
	cityName := "Тестовый Город " + testUUID
	cityID := 5

	// Создаем тестовый город
	_, err = db.Exec(ctx,
		`INSERT INTO cities (id, name) VALUES ($1, $2)`,
		cityID, cityName)
	assert.NoError(t, err)

	// Тест 1: Успешное создание PVZ
	t.Run("Success", func(t *testing.T) {
		pvz := &entity.PVZ{
			UUID:      testUUID,
			City:      entity.City{Id: cityID, Name: cityName},
			CreatedAt: time.Now(),
		}

		createdPVZ, err := repo.CreatePVZ(ctx, pvz)
		assert.NoError(t, err)
		assert.NotNil(t, createdPVZ)
		assert.Equal(t, pvz.UUID, createdPVZ.UUID)
		assert.Equal(t, pvz.City.Name, createdPVZ.City.Name)
		assert.NotZero(t, createdPVZ.CreatedAt)
	})

	// Тест 2: Ошибка при создании PVZ (если PVZ уже существует)
	t.Run("ErrorPVZAlreadyExists", func(t *testing.T) {
		pvz := &entity.PVZ{
			UUID:      testUUID,
			City:      entity.City{Id: cityID, Name: cityName},
			CreatedAt: time.Now(),
		}

		createdPVZ, err := repo.CreatePVZ(ctx, pvz)
		assert.Error(t, err)
		assert.Nil(t, createdPVZ)
		assert.Equal(t, myerrors.ErrPVZAlreadyExists, err)
	})

	// Очищаем тестовые данные
	_, err = db.Exec(ctx, `DELETE FROM pvz WHERE id = $1`, testUUID)
	assert.NoError(t, err)
	_, err = db.Exec(ctx, `DELETE FROM cities WHERE id = $1`, cityID)
	assert.NoError(t, err)
}
