package product

import (
	"avito-pvz/internal/config"
	"avito-pvz/pkg/postgres"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddProduct_Success(t *testing.T) {
	city := "1"
	pvzID := uuid.New().String()
	receivingID := uuid.New().String()
	categoryID := 1

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
	_, err = db.Exec(ctx, `INSERT INTO receivings (id, pvz_id, status, start_time) VALUES ($1, $2, $3, $4)`, receivingID, pvzID, "in_progress", "2023-01-01 00:00:00")
	assert.NoError(t, err)

	// 3. Теперь добавляем продукт
	product, err := repo.AddProduct(ctx, receivingID, categoryID)
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, receivingID, product.ReceptionID)

	// Проверяем, что продукт действительно добавлен в базу данных
	var count int
	err = db.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE receiving_id = $1`, receivingID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Удаляем продукт после теста
	_, err = db.Exec(ctx, `DELETE FROM products WHERE receiving_id = $1`, receivingID)
	assert.NoError(t, err)

	// Удаляем запись из receivings после теста
	_, err = db.Exec(ctx, `DELETE FROM receivings WHERE id = $1`, receivingID)
	assert.NoError(t, err)

	// Удаляем запись из pvz после теста
	_, err = db.Exec(ctx, `DELETE FROM pvz WHERE id = $1`, pvzID)
	assert.NoError(t, err)
}

func TestAddProduct_InvalidData(t *testing.T) {
	// Здесь вы можете протестировать некорректные данные, например, пустой receivingId или недопустимый categoryId

	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	ctx := context.Background()

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	// Пытаемся добавить продукт с пустым receivingId
	product, err := repo.AddProduct(ctx, "", 1)
	assert.Error(t, err)
	assert.Nil(t, product)

	// Пытаемся добавить продукт с недопустимым categoryId (например, отрицательное значение)
	product, err = repo.AddProduct(ctx, "receiving-123", -1)
	assert.Error(t, err)
	assert.Nil(t, product)
}
