package middlewares

import (
	"avito-pvz/internal/constants"
	myerrors "avito-pvz/internal/constants/errors"
	"avito-pvz/internal/transport/http/handlers/auth"
	"avito-pvz/pkg/logger"

	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type tokenData struct {
	jwt.RegisteredClaims        // техническое поле для пирсинга
	UserId               int    `json:"id"`
	Role                 string `json:"role"`
	CreatedAt            int64  `json:"iat"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	secretKeyString := os.Getenv("AUTH_SECRET_KEY")

	secretKey := []byte(secretKeyString)

	if secretKey == nil {
		log.Fatal("AUTH_SECRET_KEY not founded")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenHeader := r.Header.Get("Authorization") // получение данных из заголовка

		if len(accessTokenHeader) == 0 || !(strings.HasPrefix(accessTokenHeader, "Bearer ")) { // проверка, что токен начинается с корректного обозначения типа
			log.Printf("Could not get token %s", accessTokenHeader)
			http.Error(w, "Некорректный jwt", http.StatusBadRequest)
			return
		}

		accessTokenString := accessTokenHeader[7:] // извлечение самой строки токена
		token, err := jwt.ParseWithClaims(accessTokenString, &tokenData{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		ctx := r.Context()

		if data, ok := token.Claims.(*tokenData); ok && token.Valid {
			// Подготовка общих полей для логов
			logFields := []zap.Field{
				zap.Any("user_id", data.UserId),
				zap.Any("user_role", data.Role),
				zap.Any("Сейчас",time.Now().Unix()),
				zap.Any("Срок действия",time.Unix(data.CreatedAt, 0).Add(constants.ExpirationTime).Unix()),
			}

			// Проверка роли
			if data.Role != constants.Moderator {
				logger.GetLoggerFromCtx(ctx).Info(ctx, "Access denied: insufficient privileges", logFields...)

				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(auth.ErrorResponse{
					Message: myerrors.ErrAccessDenied.Error(),
				})
				return
			}

			// Проверка срока действия
			expirationTime := time.Unix(data.CreatedAt, 0).Add(constants.ExpirationTime).Unix()

			if time.Now().Unix() > expirationTime {
				logger.GetLoggerFromCtx(ctx).Info(ctx, "Access token expired", logFields...)

				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(auth.ErrorResponse{
					Message: myerrors.ErrInvalidOrExpiredJWT.Error(),
				})
				return
			}

			// Если все проверки пройдены
			ctx = context.WithValue(ctx, "userId", data.UserId)
		} else {
			// Невалидный токен
			logger.GetLoggerFromCtx(ctx).Info(ctx, "Invalid token",
				zap.Any("token_error", err.Error()))

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(auth.ErrorResponse{
				Message: myerrors.ErrInvalidOrExpiredJWT.Error(),
			})
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
