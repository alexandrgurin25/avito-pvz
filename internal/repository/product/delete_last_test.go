package product

import (
	"avito-pvz/internal/config"
	"avito-pvz/pkg/postgres"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteLastProduct_Success(t *testing.T) {
	city := "1"
	pvzID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	receivingID := "5fa85f64-5717-4562-b3fc-2c963f66afa6"
	categoryID := 1
	productID := "7fa85f64-5717-4562-b3fc-2c963f66afa6"

	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	ctx := context.Background()

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)
	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	// 1. Сначала добавляем запись в таблицу pvz
	_, err = db.Exec(ctx, `INSERT INTO pvz (id, city_id) VALUES ($1, $2)`, pvzID, city)
	assert.NoError(t, err)

	// 2. Затем добавляем запись в таблицу receivings
	_, err = db.Exec(ctx, `INSERT INTO receivings (id, pvz_id, status, start_time) VALUES ($1, $2, $3, $4)`,
		receivingID, pvzID, "in_progress", "2023-01-01 00:00:00")
	assert.NoError(t, err)

	// 3. Добавляем продукт для тестирования удаления
	_, err = db.Exec(ctx, `INSERT INTO products (id, receiving_id, category_id) VALUES ($1, $2, $3)`,
		productID, receivingID, categoryID)
	assert.NoError(t, err)

	// Проверяем, что продукт существует перед удалением
	var countBefore int
	err = db.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE id = $1`, productID).Scan(&countBefore)
	assert.NoError(t, err)
	assert.Equal(t, 1, countBefore)

	// 4. Тестируем удаление продукта
	err = repo.DeleteLastProduct(ctx, productID)
	assert.NoError(t, err)

	// Проверяем, что продукт удален
	var countAfter int
	err = db.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE id = $1`, productID).Scan(&countAfter)
	assert.NoError(t, err)
	assert.Equal(t, 0, countAfter)

	// Очищаем тестовые данные
	_, err = db.Exec(ctx, `DELETE FROM receivings WHERE id = $1`, receivingID)
	assert.NoError(t, err)

	_, err = db.Exec(ctx, `DELETE FROM pvz WHERE id = $1`, pvzID)
	assert.NoError(t, err)
}

func TestDeleteLastProduct_InvalidData(t *testing.T) {
	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	ctx := context.Background()

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	// Пытаемся удалить продукт с пустым ID
	err = repo.DeleteLastProduct(ctx, "")
	assert.Error(t, err)

	// Пытаемся удалить несуществующий продукт
	err = repo.DeleteLastProduct(ctx, "non-existent-id")
	assert.Error(t, err)
}
