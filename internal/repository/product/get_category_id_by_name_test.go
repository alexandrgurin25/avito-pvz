package product

import (
	"avito-pvz/internal/config"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/pkg/postgres"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIdCategoryByName_Success(t *testing.T) {
	tests := []struct {
		name         string
		categoryName string
		expectedId   int
	}{
		{
			name:         "Обувь",
			categoryName: "обувь",
			expectedId:   3,
		},
		{
			name:         "Одежда",
			categoryName: "одежда",
			expectedId:   2,
		},
		{
			name:         "Электроника",
			categoryName: "электроника",
			expectedId:   1,
		},
	}

	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	ctx := context.Background()

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.GetIdCategoryByName(ctx, tt.categoryName)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedId, id)
		})
	}

}

func TestGetIdCategoryByName_NotFound(t *testing.T) {
	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	ctx := context.Background()

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	// Пытаемся получить несуществующую категорию
	nonExistentName := "non_existent_category"
	id, err := repo.GetIdCategoryByName(ctx, nonExistentName)
	assert.Error(t, err)
	assert.Equal(t, myerrors.ErrInvalidProductType, err)
	assert.Equal(t, 0, id)
}

func TestGetIdCategoryByName_EmptyName(t *testing.T) {
	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	ctx := context.Background()

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	// Пытаемся получить категорию с пустым именем
	id, err := repo.GetIdCategoryByName(ctx, "")
	assert.Error(t, err)
	assert.Equal(t, 0, id)
}
