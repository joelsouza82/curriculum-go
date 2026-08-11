package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	"github.com/joelsouza82/curriculum-go/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestPersonalHandler_GetPersonals(t *testing.T) {
	mockService := new(mocks.PersonalServiceMock)
	handler := NewPersonalHandler(mockService)

	router := setupRouter()
	router.GET("/personal", handler.GetPersonals)

	birthDate := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	personalMock1 := domain.Personal{
		ID:        1,
		Name:      "User One",
		Email:     "user1@example.com",
		BirthDate: birthDate,
		LoginId:   1,
	}

	personalMock2 := domain.Personal{
		ID:        2,
		Name:      "User Two",
		Email:     "user2@example.com",
		BirthDate: birthDate,
		LoginId:   2,
	}

	t.Run("Success", func(t *testing.T) {
		mockService.On("GetPersonals").Return([]domain.Personal{personalMock1, personalMock2}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/personal", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responsePersonals []domain.Personal
		err := json.Unmarshal(w.Body.Bytes(), &responsePersonals)
		assert.NoError(t, err)
		assert.Len(t, responsePersonals, 2)
		assert.Equal(t, personalMock1.ID, responsePersonals[0].ID)
		assert.Equal(t, personalMock2.ID, responsePersonals[1].ID)

		mockService.AssertExpectations(t)
	})

	t.Run("Empty List", func(t *testing.T) {
		mockService.On("GetPersonals").Return([]domain.Personal{}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/personal", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responsePersonals []domain.Personal
		err := json.Unmarshal(w.Body.Bytes(), &responsePersonals)
		assert.NoError(t, err)
		assert.Len(t, responsePersonals, 0)

		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		internalError := errors.New("some internal error")
		mockService.On("GetPersonals").Return([]domain.Personal{}, internalError).Once()

		req, _ := http.NewRequest(http.MethodGet, "/personal", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), internalError.Error())

		mockService.AssertExpectations(t)
	})
}

func TestPersonalHandler_GetPersonalByID(t *testing.T) {
	mockService := new(mocks.PersonalServiceMock)
	handler := NewPersonalHandler(mockService)

	router := setupRouter()
	router.GET("/personal/:personalId", handler.GetPersonalByID)

	personalMock := domain.Personal{
		ID:    1,
		Email: "test.user@example.com",
	}

	t.Run("Success", func(t *testing.T) {
		mockService.On("GetPersonalByID", 1).Return(personalMock, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/personal/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responsePersonal domain.Personal
		err := json.Unmarshal(w.Body.Bytes(), &responsePersonal)
		assert.NoError(t, err)
		assert.Equal(t, personalMock, responsePersonal)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/personal/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "ID de personal inválido")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService.On("GetPersonalByID", 1).Return(domain.Personal{}, domain.ErrPersonalNotFound).Once()

		req, _ := http.NewRequest(http.MethodGet, "/personal/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Registro de personal não encontrado")

		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		internalError := errors.New("some internal error")
		mockService.On("GetPersonalByID", 1).Return(domain.Personal{}, internalError).Once()

		req, _ := http.NewRequest(http.MethodGet, "/personal/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), internalError.Error())

		mockService.AssertExpectations(t)
	})
}

func TestPersonalHandler_UpdatePersonal(t *testing.T) {
	mockService := new(mocks.PersonalServiceMock)
	handler := NewPersonalHandler(mockService)

	router := setupRouter()
	router.PUT("/personal/:personalId", handler.UpdatePersonal)

	personalToUpdate := domain.Personal{
		ID:    1,
		Email: "updated.user@example.com",
	}

	t.Run("Success", func(t *testing.T) {
		mockService.On("UpdatePersonal", mock.AnythingOfType("domain.Personal")).Return(personalToUpdate, nil).Once()

		body, _ := json.Marshal(personalToUpdate)
		req, _ := http.NewRequest(http.MethodPut, "/personal/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responsePersonal domain.Personal
		err := json.Unmarshal(w.Body.Bytes(), &responsePersonal)
		assert.NoError(t, err)
		assert.Equal(t, personalToUpdate.Email, responsePersonal.Email)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/personal/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "ID de personal inválido")
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/personal/1", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService.On("UpdatePersonal", mock.AnythingOfType("domain.Personal")).Return(domain.Personal{}, domain.ErrPersonalNotFound).Once()

		body, _ := json.Marshal(personalToUpdate)
		req, _ := http.NewRequest(http.MethodPut, "/personal/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Registro de personal não encontrado")

		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		internalError := errors.New("some internal error")
		mockService.On("UpdatePersonal", mock.AnythingOfType("domain.Personal")).Return(domain.Personal{}, internalError).Once()

		body, _ := json.Marshal(personalToUpdate)
		req, _ := http.NewRequest(http.MethodPut, "/personal/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), internalError.Error())

		mockService.AssertExpectations(t)
	})
}

func TestPersonalHandler_CreatePersonal(t *testing.T) {
	mockService := new(mocks.PersonalServiceMock)
	handler := NewPersonalHandler(mockService)

	router := setupRouter()
	router.POST("/personal", handler.CreatePersonal)

	birthDate := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	personalToCreate := domain.Personal{
		Name:      "New User",
		Email:     "newuser@example.com",
		BirthDate: birthDate,
		LoginId:   1,
	}

	t.Run("Success", func(t *testing.T) {
		createdPersonal := domain.Personal{
			ID:        10,
			Name:      "New User",
			Email:     "newuser@example.com",
			BirthDate: birthDate,
			LoginId:   1,
		}
		mockService.On("CreatePersonal", mock.AnythingOfType("domain.Personal")).Return(createdPersonal, nil).Once()

		body, _ := json.Marshal(personalToCreate)
		req, _ := http.NewRequest(http.MethodPost, "/personal", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response domain.Personal
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 10, response.ID)
		assert.Equal(t, "New User", response.Name)

		mockService.AssertExpectations(t)
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/personal", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		internalError := errors.New("some internal error")
		mockService.On("CreatePersonal", mock.AnythingOfType("domain.Personal")).Return(domain.Personal{}, internalError).Once()

		body, _ := json.Marshal(personalToCreate)
		req, _ := http.NewRequest(http.MethodPost, "/personal", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), internalError.Error())

		mockService.AssertExpectations(t)
	})
}

func TestPersonalHandler_DeletePersonal(t *testing.T) {
	mockService := new(mocks.PersonalServiceMock)
	handler := NewPersonalHandler(mockService)

	router := setupRouter()
	router.DELETE("/personal/:personalId", handler.DeletePersonal)

	personalID := 1

	t.Run("Success", func(t *testing.T) {
		mockService.On("DeletePersonal", personalID).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/personal/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/personal/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "ID de personal inválido")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService.On("DeletePersonal", personalID).Return(domain.ErrPersonalNotFound).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/personal/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Registro de personal não encontrado")
		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		internalError := errors.New("some internal error")
		mockService.On("DeletePersonal", personalID).Return(internalError).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/personal/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), internalError.Error())
		mockService.AssertExpectations(t)
	})
}
