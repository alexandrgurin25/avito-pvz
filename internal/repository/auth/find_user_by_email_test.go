package auth

import (
	"avito-pvz/internal/config"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/pkg/postgres"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindUserByEmail_Success(t *testing.T) {
	email := "test@example.com"
	passwordHash := "hashedpassword"

	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	ctx := context.Background()

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	// Сначала создаем пользователя
	_, err = repo.CreateUser(ctx, email, passwordHash, "employee")
	assert.NoError(t, err)

	// Теперь пытаемся найти пользователя по email
	user, err := repo.FindUserByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, passwordHash, user.PasswordHash)

	// Удаляем пользователя после теста
	_, err = db.Exec(ctx, `DELETE FROM users WHERE email = $1`, email)
	assert.NoError(t, err)
}

func TestFindUserByEmail_NotFound(t *testing.T) {
	email := "nonexistent@example.com"

	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	ctx := context.Background()

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	// Пытаемся найти пользователя, который не существует
	user, err := repo.FindUserByEmail(ctx, email)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, myerrors.ErrUserNotFound, err) // Проверяем, что ошибка соответствует ожидаемой
}
