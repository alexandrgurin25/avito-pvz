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

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	email := "test@example.com"
	password := "correct_password"
	userID := 1

	passwordHash, err := createHashPassword("correct_password")
	require.NoError(t, err)
	mockRepo := mocks.NewMockRepository(ctrl)
	mockUser := &entity.User{
		Id:           userID,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "user",
	}

	// Expectations
	mockRepo.EXPECT().
		FindUserByEmail(ctx, email).
		Return(mockUser, nil)

	// Mock the JWT generation (you might need to extract this to a mockable interface)
	service := NewService(mockRepo)

	// Call method
	token, err := service.Login(ctx, email, password)

	// Verify
	require.NoError(t, err)
	assert.NotNil(t, token)
}

func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	email := "nonexistent@example.com"
	password := "any_password"

	mockRepo := mocks.NewMockRepository(ctrl)

	// Expectations
	mockRepo.EXPECT().
		FindUserByEmail(ctx, email).
		Return(nil, myerrors.ErrUserNotFound)

	service := NewService(mockRepo)

	// Call method
	token, err := service.Login(ctx, email, password)

	// Verify
	require.Error(t, err)
	assert.Empty(t, token)
	assert.ErrorIs(t, err, myerrors.ErrUserNotFound)
}

func TestLogin_FindUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	email := "test@example.com"
	password := "any_password"
	dbErr := errors.New("database error")

	mockRepo := mocks.NewMockRepository(ctrl)

	// Expectations
	mockRepo.EXPECT().
		FindUserByEmail(ctx, email).
		Return(nil, dbErr)

	service := NewService(mockRepo)

	// Call method
	token, err := service.Login(ctx, email, password)

	// Verify
	require.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "failed to find user by email")
}

func TestLogin_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	email := "test@example.com"
	password := "wrong_password"
	userID := 2

	mockRepo := mocks.NewMockRepository(ctrl)
	passwordHash, err := createHashPassword("correct_password")
	require.NoError(t, err)

	mockUser := &entity.User{
		Id:           userID,
		Email:        email,
		PasswordHash: passwordHash, // Assume this is hashed
		Role:         "user",
	}

	// Expectations
	mockRepo.EXPECT().
		FindUserByEmail(ctx, email).
		Return(mockUser, nil)

	service := NewService(mockRepo)

	// Call method
	token, err := service.Login(ctx, email, password)
	require.Error(t, err)

	// Verify
	require.Error(t, err)
	assert.Empty(t, token)
	assert.ErrorIs(t, err, myerrors.ErrWrongPassword)
}
