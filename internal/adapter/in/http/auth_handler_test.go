package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	"github.com/joelsouza82/curriculum-go/internal/mocks"

	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login(t *testing.T) {
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)

	router := setupRouter()
	router.POST("/auth/login", handler.Login)

	credentials := authRequest{
		Email:    "user@example.com",
		Password: "correct-password",
	}

	t.Run("Success", func(t *testing.T) {
		mockService.On("Authenticate", credentials.Email, credentials.Password).Return("signed-jwt-token", nil).Once()

		body, _ := json.Marshal(credentials)
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "signed-jwt-token", response["token"])

		mockService.AssertExpectations(t)
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		mockService.On("Authenticate", credentials.Email, credentials.Password).Return("", domain.ErrInvalidCredentials).Once()

		body, _ := json.Marshal(credentials)
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "credenciais inválidas")

		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		internalError := errors.New("database error")
		mockService.On("Authenticate", credentials.Email, credentials.Password).Return("", internalError).Once()

		body, _ := json.Marshal(credentials)
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), internalError.Error())

		mockService.AssertExpectations(t)
	})
}
