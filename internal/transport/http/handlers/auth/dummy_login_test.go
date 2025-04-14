package auth_test

import (
	"avito-pvz/internal/constants"
	"errors"

	service "avito-pvz/internal/service/auth/mocks"
	dto "avito-pvz/internal/transport/http/dto/auth"
	message "avito-pvz/internal/transport/http/dto/error"
	"avito-pvz/internal/transport/http/handlers/auth"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDummyLogin(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      dto.DummyLoginRequest
		mockReturnToken  string
		mockReturnError  error
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:             "Успешный логин для Employee",
			requestBody:      dto.DummyLoginRequest{Role: constants.Employee},
			mockReturnToken:  "employee_token",
			mockReturnError:  nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: dto.DummyLoginResponse{Token: "employee_token"},
		},
		{
			name:             "Успешный логин для Moderator",
			requestBody:      dto.DummyLoginRequest{Role: constants.Moderator},
			mockReturnToken:  "moderator_token",
			mockReturnError:  nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: dto.DummyLoginResponse{Token: "moderator_token"},
		},

		{
			name:             "Рандомная роль",
			requestBody:      dto.DummyLoginRequest{Role: "Random"},
			mockReturnToken:  "",
			mockReturnError:  nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: message.ErrorResponse{Message: "недопустимая роль"},
		},
		{
			name:             "Ошибка при создании логина",
			requestBody:      dto.DummyLoginRequest{Role: constants.Employee},
			mockReturnToken:  "",
			mockReturnError:  errors.New("внутренняя ошибка сервиса"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: message.ErrorResponse{Message: "внутренняя ошибка сервиса"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := service.NewMockService(ctrl)

			// Only expect the mock call if we're not testing the invalid role case
			if tt.expectedStatus == http.StatusOK || tt.name == "Ошибка при создании логина" {
				mockService.EXPECT().CreateDummyLogin(tt.requestBody.Role).Return(tt.mockReturnToken, tt.mockReturnError)
			}

			handler := auth.NewHandler(mockService)
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.DummyLogin(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			// Handle response decoding based on expected type
			if res.StatusCode == http.StatusOK {
				var response dto.DummyLoginResponse
				err := json.NewDecoder(res.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			} else {
				var response message.ErrorResponse
				err := json.NewDecoder(res.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}
