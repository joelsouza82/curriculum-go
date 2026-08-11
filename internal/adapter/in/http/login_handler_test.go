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
	"github.com/stretchr/testify/mock"
)

func TestLoginHandler_GetLogins(t *testing.T) {
	mockService := new(mocks.LoginServiceMock)
	handler := NewLoginHandler(mockService)

	router := setupRouter()
	router.GET("/login", handler.GetLogins)

	loginMock1 := domain.Login{
		ID:       1,
		Email:    "user1@example.com",
		Password: "pass1",
	}

	loginMock2 := domain.Login{
		ID:       2,
		Email:    "user2@example.com",
		Password: "pass2",
	}

	t.Run("Success", func(t *testing.T) {
		mockService.On("GetLogins").Return([]domain.Login{loginMock1, loginMock2}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/login", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseLogins []domain.Login
		err := json.Unmarshal(w.Body.Bytes(), &responseLogins)
		assert.NoError(t, err)
		assert.Len(t, responseLogins, 2)
		assert.Equal(t, loginMock1.ID, responseLogins[0].ID)
		assert.Equal(t, loginMock2.ID, responseLogins[1].ID)

		mockService.AssertExpectations(t)
	})

	t.Run("Empty List", func(t *testing.T) {
		mockService.On("GetLogins").Return([]domain.Login{}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/login", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseLogins []domain.Login
		err := json.Unmarshal(w.Body.Bytes(), &responseLogins)
		assert.NoError(t, err)
		assert.Len(t, responseLogins, 0)

		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		internalError := errors.New("database error")
		mockService.On("GetLogins").Return([]domain.Login{}, internalError).Once()

		req, _ := http.NewRequest(http.MethodGet, "/login", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), internalError.Error())

		mockService.AssertExpectations(t)
	})
}

func TestLoginHandler_GetLoginByID(t *testing.T) {
	mockService := new(mocks.LoginServiceMock)
	handler := NewLoginHandler(mockService)

	router := setupRouter()
	router.GET("/login/:loginId", handler.GetLoginByID)

	loginMock := domain.Login{
		ID:       1,
		Email:    "test@example.com",
		Password: "pass123",
	}

	t.Run("Success", func(t *testing.T) {
		mockService.On("GetLoginByID", 1).Return(loginMock, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/login/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseLogin domain.Login
		err := json.Unmarshal(w.Body.Bytes(), &responseLogin)
		assert.NoError(t, err)
		assert.Equal(t, loginMock, responseLogin)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/login/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "ID de login inválido")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService.On("GetLoginByID", 1).Return(domain.Login{}, domain.ErrLoginNotFound).Once()

		req, _ := http.NewRequest(http.MethodGet, "/login/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Registro de login não encontrado")

		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		internalError := errors.New("database error")
		mockService.On("GetLoginByID", 1).Return(domain.Login{}, internalError).Once()

		req, _ := http.NewRequest(http.MethodGet, "/login/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), internalError.Error())

		mockService.AssertExpectations(t)
	})
}

func TestLoginHandler_CreateLogin(t *testing.T) {
	mockService := new(mocks.LoginServiceMock)
	handler := NewLoginHandler(mockService)

	router := setupRouter()
	router.POST("/login", handler.CreateLogin)

	loginToCreate := domain.Login{
		Email:    "newuser@example.com",
		Password: "newpass123",
	}

	t.Run("Success", func(t *testing.T) {
		createdLogin := domain.Login{
			ID:       10,
			Email:    "newuser@example.com",
			Password: "newpass123",
		}
		mockService.On("CreateLogin", mock.AnythingOfType("domain.Login")).Return(createdLogin, nil).Once()

		body, _ := json.Marshal(loginToCreate)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response domain.Login
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 10, response.ID)
		assert.Equal(t, "newuser@example.com", response.Email)

		mockService.AssertExpectations(t)
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		internalError := errors.New("database error")
		mockService.On("CreateLogin", mock.AnythingOfType("domain.Login")).Return(domain.Login{}, internalError).Once()

		body, _ := json.Marshal(loginToCreate)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), internalError.Error())

		mockService.AssertExpectations(t)
	})
}

func TestLoginHandler_UpdateLogin(t *testing.T) {
	mockService := new(mocks.LoginServiceMock)
	handler := NewLoginHandler(mockService)

	router := setupRouter()
	router.PUT("/login/:loginId", handler.UpdateLogin)

	loginToUpdate := domain.Login{
		ID:       1,
		Email:    "updated@example.com",
		Password: "updatedpass",
	}

	t.Run("Success", func(t *testing.T) {
		mockService.On("UpdateLogin", mock.AnythingOfType("domain.Login")).Return(loginToUpdate, nil).Once()

		body, _ := json.Marshal(loginToUpdate)
		req, _ := http.NewRequest(http.MethodPut, "/login/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseLogin domain.Login
		err := json.Unmarshal(w.Body.Bytes(), &responseLogin)
		assert.NoError(t, err)
		assert.Equal(t, loginToUpdate.Email, responseLogin.Email)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/login/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "ID de login inválido")
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/login/1", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService.On("UpdateLogin", mock.AnythingOfType("domain.Login")).Return(domain.Login{}, domain.ErrLoginNotFound).Once()

		body, _ := json.Marshal(loginToUpdate)
		req, _ := http.NewRequest(http.MethodPut, "/login/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Registro de login não encontrado")

		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		internalError := errors.New("database error")
		mockService.On("UpdateLogin", mock.AnythingOfType("domain.Login")).Return(domain.Login{}, internalError).Once()

		body, _ := json.Marshal(loginToUpdate)
		req, _ := http.NewRequest(http.MethodPut, "/login/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), internalError.Error())

		mockService.AssertExpectations(t)
	})
}

func TestLoginHandler_DeleteLogin(t *testing.T) {
	mockService := new(mocks.LoginServiceMock)
	handler := NewLoginHandler(mockService)

	router := setupRouter()
	router.DELETE("/login/:loginId", handler.DeleteLogin)

	loginID := 1

	t.Run("Success", func(t *testing.T) {
		mockService.On("DeleteLogin", loginID).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/login/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/login/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "ID de login inválido")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService.On("DeleteLogin", loginID).Return(domain.ErrLoginNotFound).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/login/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Registro de login não encontrado")
		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		internalError := errors.New("database error")
		mockService.On("DeleteLogin", loginID).Return(internalError).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/login/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), internalError.Error())
		mockService.AssertExpectations(t)
	})
}
