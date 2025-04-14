package auth_test

import (
	myerrors "avito-pvz/internal/constants/errors"
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

func TestLogin(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      dto.LoginRequest
		mockReturnToken  string
		mockReturnError  error
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:             "Successful login",
			requestBody:      dto.LoginRequest{Email: "alex@avito.ru", Password: "1234"},
			mockReturnToken:  "login_token",
			mockReturnError:  nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: dto.LoginResponse{Token: "login_token"},
		},
		{
			name:             "Empty email and password",
			requestBody:      dto.LoginRequest{Email: "", Password: ""},
			mockReturnToken:  "",
			mockReturnError:  nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: message.ErrorResponse{Message: myerrors.ErrEmailOrPasswordEmpty.Error()},
		},
		{
			name:             "Authentication error",
			requestBody:      dto.LoginRequest{Email: "alex@avito.ru", Password: "1234"},
			mockReturnToken:  "",
			mockReturnError:  errors.New("invalid credentials"),
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: message.ErrorResponse{Message: "invalid credentials"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := service.NewMockService(ctrl)

			if tt.expectedStatus == http.StatusOK || tt.expectedStatus == http.StatusUnauthorized {
				mockService.EXPECT().
					Login(gomock.Any(), tt.requestBody.Email, tt.requestBody.Password).
					Return(tt.mockReturnToken, tt.mockReturnError)
			}

			handler := auth.NewHandler(mockService)
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			// Handle response decoding based on expected type
			if res.StatusCode == http.StatusOK {
				var response dto.LoginResponse
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
