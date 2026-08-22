package service

import (
	"errors"
	"testing"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	"github.com/joelsouza82/curriculum-go/internal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginService_GetLoginByID(t *testing.T) {
	mockRepo := new(mocks.LoginRepositoryMock)
	svc := NewLoginService(mockRepo)

	loginMock := domain.Login{
		ID:       1,
		Email:    "test@test.com",
		Password: "pass123",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetLoginByID", 1).Return(loginMock, nil).Once()

		login, err := svc.GetLoginByID(1)

		assert.NoError(t, err)
		assert.Equal(t, loginMock, login)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repoError := errors.New("login not found")
		mockRepo.On("GetLoginByID", 1).Return(domain.Login{}, repoError).Once()

		_, err := svc.GetLoginByID(1)

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestLoginService_GetLogins(t *testing.T) {
	mockRepo := new(mocks.LoginRepositoryMock)
	svc := NewLoginService(mockRepo)

	loginMock1 := domain.Login{
		ID:       1,
		Email:    "user1@test.com",
		Password: "pass1",
	}

	loginMock2 := domain.Login{
		ID:       2,
		Email:    "user2@test.com",
		Password: "pass2",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetLogins").Return([]domain.Login{loginMock1, loginMock2}, nil).Once()

		logins, err := svc.GetLogins()

		assert.NoError(t, err)
		assert.Len(t, logins, 2)
		assert.Equal(t, loginMock1, logins[0])
		assert.Equal(t, loginMock2, logins[1])
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repoError := errors.New("database error")
		mockRepo.On("GetLogins").Return([]domain.Login{}, repoError).Once()

		_, err := svc.GetLogins()

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty List", func(t *testing.T) {
		mockRepo.On("GetLogins").Return([]domain.Login{}, nil).Once()

		logins, err := svc.GetLogins()

		assert.NoError(t, err)
		assert.Len(t, logins, 0)
		mockRepo.AssertExpectations(t)
	})
}

func TestLoginService_CreateLogin(t *testing.T) {
	mockRepo := new(mocks.LoginRepositoryMock)
	svc := NewLoginService(mockRepo)

	loginToCreate := domain.Login{
		Email:    "newuser@test.com",
		Password: "newpass",
	}

	passwordMatcher := func(l domain.Login) bool {
		return l.Email == loginToCreate.Email &&
			bcrypt.CompareHashAndPassword([]byte(l.Password), []byte(loginToCreate.Password)) == nil
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetLoginByEmail", loginToCreate.Email).Return(domain.Login{}, domain.ErrLoginNotFound).Once()
		mockRepo.On("CreateLogin", mock.MatchedBy(passwordMatcher)).Return(10, nil).Once()

		createdLogin, err := svc.CreateLogin(loginToCreate)

		assert.NoError(t, err)
		assert.Equal(t, 10, createdLogin.ID)
		assert.Equal(t, loginToCreate.Email, createdLogin.Email)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(createdLogin.Password), []byte(loginToCreate.Password)))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		existingLogin := domain.Login{ID: 5, Email: loginToCreate.Email, Password: "hashed"}
		mockRepo.On("GetLoginByEmail", loginToCreate.Email).Return(existingLogin, nil).Once()

		_, err := svc.CreateLogin(loginToCreate)

		assert.ErrorIs(t, err, domain.ErrEmailAlreadyExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repoError := errors.New("database error")
		mockRepo.On("GetLoginByEmail", loginToCreate.Email).Return(domain.Login{}, domain.ErrLoginNotFound).Once()
		mockRepo.On("CreateLogin", mock.MatchedBy(passwordMatcher)).Return(0, repoError).Once()

		_, err := svc.CreateLogin(loginToCreate)

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Email Lookup Error", func(t *testing.T) {
		repoError := errors.New("database error")
		mockRepo.On("GetLoginByEmail", loginToCreate.Email).Return(domain.Login{}, repoError).Once()

		_, err := svc.CreateLogin(loginToCreate)

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestLoginService_UpdateLogin(t *testing.T) {
	mockRepo := new(mocks.LoginRepositoryMock)
	svc := NewLoginService(mockRepo)

	loginToUpdate := domain.Login{
		ID:       1,
		Email:    "updated@test.com",
		Password: "updatedpass",
	}

	updateMatcher := func(l domain.Login) bool {
		return l.ID == loginToUpdate.ID && l.Email == loginToUpdate.Email &&
			bcrypt.CompareHashAndPassword([]byte(l.Password), []byte(loginToUpdate.Password)) == nil
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("UpdateLogin", mock.MatchedBy(updateMatcher)).Return(loginToUpdate, nil).Once()

		updatedLogin, err := svc.UpdateLogin(loginToUpdate)

		assert.NoError(t, err)
		assert.Equal(t, loginToUpdate.ID, updatedLogin.ID)
		assert.Equal(t, loginToUpdate.Email, updatedLogin.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repoError := errors.New("login not found")
		mockRepo.On("UpdateLogin", mock.MatchedBy(updateMatcher)).Return(domain.Login{}, repoError).Once()

		_, err := svc.UpdateLogin(loginToUpdate)

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestLoginService_DeleteLogin(t *testing.T) {
	mockRepo := new(mocks.LoginRepositoryMock)
	svc := NewLoginService(mockRepo)
	loginID := 1

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("DeleteLogin", loginID).Return(nil).Once()

		err := svc.DeleteLogin(loginID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		repoError := errors.New("login not found")
		mockRepo.On("DeleteLogin", loginID).Return(repoError).Once()

		err := svc.DeleteLogin(loginID)

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		mockRepo.AssertExpectations(t)
	})
}
