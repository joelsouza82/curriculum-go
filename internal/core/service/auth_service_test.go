package service

import (
	"errors"
	"testing"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	"github.com/joelsouza82/curriculum-go/internal/mocks"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Authenticate(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	loginMock := domain.Login{
		ID:       1,
		Email:    "user@example.com",
		Password: string(hashedPassword),
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mocks.LoginRepositoryMock)
		svc := NewAuthService(mockRepo, "test-secret")

		mockRepo.On("GetLoginByEmail", loginMock.Email).Return(loginMock, nil).Once()

		token, err := svc.Authenticate(loginMock.Email, "correct-password")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte("test-secret"), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsed.Valid)

		claims := parsed.Claims.(jwt.MapClaims)
		assert.Equal(t, loginMock.Email, claims["email"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		mockRepo := new(mocks.LoginRepositoryMock)
		svc := NewAuthService(mockRepo, "test-secret")

		mockRepo.On("GetLoginByEmail", loginMock.Email).Return(loginMock, nil).Once()

		token, err := svc.Authenticate(loginMock.Email, "wrong-password")

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo := new(mocks.LoginRepositoryMock)
		svc := NewAuthService(mockRepo, "test-secret")

		mockRepo.On("GetLoginByEmail", "missing@example.com").Return(domain.Login{}, domain.ErrLoginNotFound).Once()

		token, err := svc.Authenticate("missing@example.com", "any-password")

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(mocks.LoginRepositoryMock)
		svc := NewAuthService(mockRepo, "test-secret")

		repoError := errors.New("database error")
		mockRepo.On("GetLoginByEmail", loginMock.Email).Return(domain.Login{}, repoError).Once()

		token, err := svc.Authenticate(loginMock.Email, "correct-password")

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})
}
