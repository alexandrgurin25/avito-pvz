package auth

import (
	"avito-pvz/internal/config"
	"avito-pvz/pkg/postgres"
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	email := "test@example.com"
	passwordHash := "hashedpassword"
	role := "employee"

	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	ctx := context.Background()

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	user, err := repo.CreateUser(ctx, email, passwordHash, role)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}
	if user.Role != role {
		t.Errorf("Expected role %s, got %s", role, user.Role)
	}

	// Проверка, что пользователь действительно добавлен в базу данных
	var count int
	err = db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE email = $1`, email).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query user count: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected user count 1, got %d", count)
	}

	_, err = db.Exec(ctx, `DELETE FROM users WHERE email = $1`, email)
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}

}

func TestCreateUser_AlreadyExists(t *testing.T) {
	email := "test@example.com"
	passwordHash := "hashedpassword"
	role := "employee"

	pathToEnv := "../../../config/test.env"
	pathToMigrate := "file://../../../db/migrations"

	ctx := context.Background()

	cfg, err := config.NewTest(pathToEnv)
	assert.NoError(t, err)

	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	repo := NewRepository(db)

	// Сначала создаем пользователя
	_, err = repo.CreateUser(ctx, email, passwordHash, role)
	assert.NoError(t, err)

	// Теперь пытаемся создать пользователя с тем же email
	user, err := repo.CreateUser(ctx, email, passwordHash, role)
	assert.Error(t, err) // Ожидаем ошибку
	assert.Nil(t, user)  // Пользователь не должен быть создан

	// Проверка, что ошибка связана с уникальностью email
	if err != nil {
		assert.Contains(t, err.Error(), "unique constraint") // Замените на конкретное сообщение об ошибке, если необходимо
	}

	// Удаляем пользователя после теста
	_, err = db.Exec(ctx, `DELETE FROM users WHERE email = $1`, email)
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
}
