package auth_test

import (
	service "avito-pvz/internal/service/auth/mocks"
	dto "avito-pvz/internal/transport/http/dto/auth"
	message "avito-pvz/internal/transport/http/dto/error"
	"avito-pvz/internal/transport/http/handlers/auth"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      dto.RegisterRequest
		mockReturnError  error
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name: "Successful registration",
			requestBody: dto.RegisterRequest{
				Email:    "alex@avito.ru",
				Password: "1234",
				Role:     "employee",
			},
			mockReturnError:  nil,
			expectedStatus:   http.StatusCreated,
			expectedResponse: nil,
		},

		{
			name: "Empty email and password",
			requestBody: dto.RegisterRequest{
				Email:    "",
				Password: "",
			},
			mockReturnError:  nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: message.ErrorResponse{Message: "email пользователя и пароль обязательны"},
		},
		{
			name: "Неверная роль",
			requestBody: dto.RegisterRequest{
				Email:    "alex@avito.ru",
				Password: "1234",
				Role:     "Продавец",
			},
			mockReturnError:  errors.New("registration failed"),
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: message.ErrorResponse{Message: "недопустимая роль"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := service.NewMockService(ctrl)

			if tt.expectedStatus == http.StatusCreated || tt.expectedStatus == http.StatusInternalServerError {
				mockService.EXPECT().
					Register(gomock.Any(), tt.requestBody.Email, tt.requestBody.Password, tt.requestBody.Role).
					Return(tt.mockReturnError)
			}

			handler := auth.NewHandler(mockService)
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.Register(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			// Обработка декодирования ответа в зависимости от ожидаемого типа
			if res.StatusCode == http.StatusCreated {
			} else {
				var response message.ErrorResponse
				err := json.NewDecoder(res.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}
