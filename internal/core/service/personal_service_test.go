package service

import (
	"errors"
	"testing"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	"github.com/joelsouza82/curriculum-go/internal/mocks"

	"github.com/stretchr/testify/assert"
)

func TestPersonalService_GetPersonalByID(t *testing.T) {
	mockRepo := new(mocks.PersonalRepositoryMock)
	svc := NewPersonalService(mockRepo)

	personalMock := domain.Personal{
		ID:    1,
		Email: "test@test.com",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetPersonalByID", 1).Return(personalMock, nil).Once()

		personal, err := svc.GetPersonalByID(1)

		assert.NoError(t, err)
		assert.Equal(t, personalMock, personal)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repoError := errors.New("record not found")
		mockRepo.On("GetPersonalByID", 1).Return(domain.Personal{}, repoError).Once()

		_, err := svc.GetPersonalByID(1)

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPersonalService_UpdatePersonal(t *testing.T) {
	mockRepo := new(mocks.PersonalRepositoryMock)
	svc := NewPersonalService(mockRepo)

	personalToUpdate := domain.Personal{
		ID:    1,
		Email: "updated@test.com",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("UpdatePersonal", personalToUpdate).Return(personalToUpdate, nil).Once()

		updatedPersonal, err := svc.UpdatePersonal(personalToUpdate)

		assert.NoError(t, err)
		assert.Equal(t, personalToUpdate.ID, updatedPersonal.ID)
		assert.Equal(t, personalToUpdate.Email, updatedPersonal.Email)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repoError := errors.New("database connection failed")

		mockRepo.On("UpdatePersonal", personalToUpdate).Return(domain.Personal{}, repoError).Once()

		_, err := svc.UpdatePersonal(personalToUpdate)

		assert.Error(t, err)
		assert.Equal(t, repoError, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestPersonalService_DeletePersonal(t *testing.T) {
	mockRepo := new(mocks.PersonalRepositoryMock)
	svc := NewPersonalService(mockRepo)
	personalID := 1

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("DeletePersonal", personalID).Return(nil).Once()

		err := svc.DeletePersonal(personalID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repoError := errors.New("record not found")
		mockRepo.On("DeletePersonal", personalID).Return(repoError).Once()

		err := svc.DeletePersonal(personalID)

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		mockRepo.AssertExpectations(t)
	})
}
