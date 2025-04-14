package auth

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/pkg/logger"
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

func (s *authService) Register(ctx context.Context, email string, password string, role string) error {
	// Поиск пользователя в бд по username
	user, err := s.authRepository.FindUserByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, myerrors.ErrUserNotFound) {
			logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to find user by email:", zap.Error(err))
			return fmt.Errorf("failed to find user by email: %v", err)
		}
		err = nil
	}

	if user != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Пользователь с такими данными уже существует")
		return myerrors.ErrUserAlreadyExists
	}

	// Проверка пароля
	passwordHash, err := createHashPassword(password)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to create hash password", zap.Error(err))
		return fmt.Errorf("failed to create hash password: %v", err)
	}

	_, err = s.authRepository.CreateUser(ctx, email, passwordHash, role)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to create user", zap.Error(err))
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}
