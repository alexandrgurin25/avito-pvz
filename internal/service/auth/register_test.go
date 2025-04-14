package auth

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/entity"
	"avito-pvz/internal/repository/auth/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	email := "test@example.com"
	password := "securePassword123"
	role := "user"

	mockRepo := mocks.NewMockRepository(ctrl)

	// 1. Ожидаем вызов FindUserByEmail с возвратом ErrUserNotFound
	mockRepo.EXPECT().
		FindUserByEmail(ctx, email).
		Return(nil, myerrors.ErrUserNotFound)

	// 2. Ожидаем вызов CreateUser с любым хешем пароля
	mockRepo.EXPECT().
		CreateUser(ctx, email, gomock.Any(), role).
		Return(nil, nil) // Успешное создание пользователя

	service := NewService(mockRepo)

	err := service.Register(ctx, email, password, role)
	require.NoError(t, err)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	email := "existing@example.com"
	password := "anyPassword"
	role := "user"

	mockRepo := mocks.NewMockRepository(ctrl)
	mockUser := &entity.User{
		Email: email,
	}

	// Expectations
	mockRepo.EXPECT().
		FindUserByEmail(ctx, email).
		Return(mockUser, nil)

	service := NewService(mockRepo)

	// Call method
	err := service.Register(ctx, email, password, role)

	// Verify
	require.Error(t, err)
	assert.ErrorIs(t, err, myerrors.ErrUserAlreadyExists)
}

func TestRegister_FindUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	email := "test@example.com"
	password := "securePassword123"
	role := "user"
	dbErr := errors.New("database error")

	mockRepo := mocks.NewMockRepository(ctrl)

	// Expectations
	mockRepo.EXPECT().
		FindUserByEmail(ctx, email).
		Return(nil, dbErr)

	service := NewService(mockRepo)

	// Call method
	err := service.Register(ctx, email, password, role)

	// Verify
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find user by email")
}
