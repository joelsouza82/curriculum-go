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
		assert.NotContains(t, w.Body.String(), "pass1")
		assert.NotContains(t, w.Body.String(), "password")

		var responseLogins []loginResponse
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

		var responseLogins []loginResponse
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
		assert.NotContains(t, w.Body.String(), "pass123")

		var responseLogin loginResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseLogin)
		assert.NoError(t, err)
		assert.Equal(t, loginMock.ID, responseLogin.ID)
		assert.Equal(t, loginMock.Email, responseLogin.Email)

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

	loginToCreate := loginRequest{
		Email:    "newuser@example.com",
		Password: "newpass123",
	}

	t.Run("Success", func(t *testing.T) {
		createdLogin := domain.Login{
			ID:       10,
			Email:    "newuser@example.com",
			Password: "$2a$10$hashedpasswordvalue",
		}
		mockService.On("CreateLogin", mock.MatchedBy(func(l domain.Login) bool {
			return l.Email == loginToCreate.Email && l.Password == loginToCreate.Password
		})).Return(createdLogin, nil).Once()

		body, _ := json.Marshal(loginToCreate)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NotContains(t, w.Body.String(), "hashedpasswordvalue")

		var response loginResponse
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
		mockService.On("CreateLogin", mock.MatchedBy(func(l domain.Login) bool {
			return l.Email == loginToCreate.Email && l.Password == loginToCreate.Password
		})).Return(domain.Login{}, internalError).Once()

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

	loginToUpdate := loginRequest{
		Email:    "updated@example.com",
		Password: "updatedpass",
	}

	t.Run("Success", func(t *testing.T) {
		updatedLogin := domain.Login{ID: 1, Email: loginToUpdate.Email, Password: "$2a$10$hashedpasswordvalue"}
		mockService.On("UpdateLogin", mock.MatchedBy(func(l domain.Login) bool {
			return l.ID == 1 && l.Email == loginToUpdate.Email && l.Password == loginToUpdate.Password
		})).Return(updatedLogin, nil).Once()

		body, _ := json.Marshal(loginToUpdate)
		req, _ := http.NewRequest(http.MethodPut, "/login/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotContains(t, w.Body.String(), "hashedpasswordvalue")

		var responseLogin loginResponse
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
		mockService.On("UpdateLogin", mock.MatchedBy(func(l domain.Login) bool {
			return l.Email == loginToUpdate.Email && l.Password == loginToUpdate.Password
		})).Return(domain.Login{}, domain.ErrLoginNotFound).Once()

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
		mockService.On("UpdateLogin", mock.MatchedBy(func(l domain.Login) bool {
			return l.Email == loginToUpdate.Email && l.Password == loginToUpdate.Password
		})).Return(domain.Login{}, internalError).Once()

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
