package pvz

import (
	"avito-pvz/internal/config"
	"avito-pvz/pkg/postgres"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetPVZsWithFilters(t *testing.T) {
	// Подготовка тестовых данных
	ctx := context.Background()
	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	// Создаем тестовые данные
	pvzID := uuid.New().String()
	cityID := 1                                  // Предполагаем, что город с ID 1 существует
	startTime := time.Now().Add(-24 * time.Hour) // 24 часа назад
	endTime := time.Now()                        // Текущий момент

	// Вставляем тестовые данные в базу
	_, err = db.Exec(ctx, `INSERT INTO pvz (id, city_id, created_at) VALUES ($1, $2, $3)`, pvzID, cityID, time.Now())
	assert.NoError(t, err)

	receivingID := uuid.New().String()
	_, err = db.Exec(ctx, `INSERT INTO receivings (id, pvz_id, start_time, end_time, status) VALUES ($1, $2, $3, $4, $5)`,
		receivingID, pvzID, startTime, endTime, "in_progress")
	assert.NoError(t, err)

	productID := uuid.New().String()
	_, err = db.Exec(ctx, `INSERT INTO products (id, receiving_id, category_id, added_at) VALUES ($1, $2, $3, $4)`,
		productID, receivingID, 1, time.Now())
	assert.NoError(t, err)

	// Тест 1: Успешное получение PVZ с фильтрами
	t.Run("Success", func(t *testing.T) {
		pvzs, err := repo.GetPVZsWithFilters(ctx, &startTime, &endTime, 1, 10)
		assert.NoError(t, err)
		assert.NotEmpty(t, pvzs)                           // Проверяем, что результат не пустой
		assert.Equal(t, pvzID, pvzs[0].UUID)               // Проверяем, что UUID совпадает
		assert.Equal(t, "Москва", pvzs[0].City.Name)       // Проверяем, что имя города совпадает
		assert.NotEmpty(t, pvzs[0].Receptions)             // Проверяем, что есть полученные данные
		assert.NotEmpty(t, pvzs[0].Receptions[0].Products) // Проверяем, что есть продукты
	})

	// Тест 2: Нет PVZ по заданным фильтрам
	t.Run("NoPVZsFound", func(t *testing.T) {
		futureStartTime := time.Now().Add(24 * time.Hour) // Время в будущем
		pvzs, err := repo.GetPVZsWithFilters(ctx, &futureStartTime, &endTime, 1, 10)
		assert.NoError(t, err)
		assert.Empty(t, pvzs) // Проверяем, что результат пустой
	})

	// Очищаем тестовые данные
	_, err = db.Exec(ctx, `DELETE FROM products WHERE id = $1`, productID)
	assert.NoError(t, err)

	_, err = db.Exec(ctx, `DELETE FROM receivings WHERE id = $1`, receivingID)
	assert.NoError(t, err)

	_, err = db.Exec(ctx, `DELETE FROM pvz WHERE id = $1`, pvzID)
	assert.NoError(t, err)
}
