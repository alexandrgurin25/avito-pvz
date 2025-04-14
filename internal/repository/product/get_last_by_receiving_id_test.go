package product

import (
	"avito-pvz/internal/config"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/pkg/postgres"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetLastProductByReceigingId(t *testing.T) {
	// Подготовка тестовых данных

	pvzID := uuid.New().String()
	receivingID := uuid.New().String()
	categoryID := "1"
	productID := uuid.New().String()

	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	ctx := context.Background()

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	_, err = db.Exec(ctx, `INSERT INTO pvz (id, city_id) VALUES ($1, $2)`, pvzID, 1)
	assert.NoError(t, err)

	_, err = db.Exec(ctx, `INSERT INTO receivings (id, pvz_id, status, start_time) VALUES ($1, $2, $3, $4)`,
		receivingID, pvzID, "in_progress", "2025-01-01 00:00:00")
	assert.NoError(t, err)

	// Тест 1: Успешное получение последнего продукта
	t.Run("Success", func(t *testing.T) {
		// Добавляем два продукта с разным временем добавления
		_, err = db.Exec(ctx, `INSERT INTO products (id, receiving_id, category_id, added_at) VALUES 
			($1, $2, $3, $4)`,
			productID, receivingID, categoryID, time.Now().Add(-1*time.Hour))
		assert.NoError(t, err)

		product, err := repo.GetLastProductByReceigingId(ctx, receivingID)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, receivingID, product.ReceptionID)
		assert.Equal(t, categoryID, product.Category)

		// Очищаем добавленные продукты
		_, err = db.Exec(ctx, `DELETE FROM products WHERE receiving_id = $1`, receivingID)
		assert.NoError(t, err)
	})

	// Тест 2: Нет продуктов для данного receiving_id
	t.Run("NoProducts", func(t *testing.T) {
		product, err := repo.GetLastProductByReceigingId(ctx, receivingID)
		assert.Error(t, err)
		assert.Equal(t, myerrors.ErrNoProductsToDelete, err)
		assert.Nil(t, product)
	})

	// Тест 3: Пустой receiving_id
	t.Run("EmptyReceivingID", func(t *testing.T) {
		product, err := repo.GetLastProductByReceigingId(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, product)
	})

	// Очищаем тестовые данные
	_, err = db.Exec(ctx, `DELETE FROM receivings WHERE id = $1`, receivingID)
	assert.NoError(t, err)

	_, err = db.Exec(ctx, `DELETE FROM pvz WHERE id = $1`, pvzID)
	assert.NoError(t, err)
}
