package auth

import (
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/pkg/logger"
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

func (s *authService) Login(ctx context.Context, email string, password string) (string, error) {
	token := ""
	// Поиск пользователя в бд по username
	user, err := s.authRepository.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserNotFound) {
			return "", err
		}
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to find user by email:", zap.Error(err))
		return "", fmt.Errorf("failed to find user by email: %v", err)
	}

	// Проверка пароля
	if checkPasswordHash(password, user.PasswordHash) {
		token, err = generateJWT(user.Id, user.Role)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to genetate jwt:", zap.Error(err))
			return "", fmt.Errorf("failed to genetate jwt:: %v", err)
		}
	} else {
		return "", myerrors.ErrWrongPassword
	}

	return token, nil
}
