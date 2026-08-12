package mocks

import (
	"github.com/joelsouza82/curriculum-go/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type LoginRepositoryMock struct {
	mock.Mock
}

func (m *LoginRepositoryMock) GetLogins() ([]domain.Login, error) {
	args := m.Called()
	if args.Get(1) != nil {
		return []domain.Login{}, args.Error(1)
	}
	return args.Get(0).([]domain.Login), nil
}

func (m *LoginRepositoryMock) GetLoginByID(id int) (domain.Login, error) {
	args := m.Called(id)
	if args.Get(1) != nil {
		return domain.Login{}, args.Error(1)
	}
	return args.Get(0).(domain.Login), nil
}

func (m *LoginRepositoryMock) GetLoginByEmail(email string) (domain.Login, error) {
	args := m.Called(email)
	if args.Get(1) != nil {
		return domain.Login{}, args.Error(1)
	}
	return args.Get(0).(domain.Login), nil
}

func (m *LoginRepositoryMock) CreateLogin(login domain.Login) (int, error) {
	args := m.Called(login)
	if args.Get(1) != nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int), nil
}

func (m *LoginRepositoryMock) UpdateLogin(login domain.Login) (domain.Login, error) {
	args := m.Called(login)
	if args.Get(1) != nil {
		return domain.Login{}, args.Error(1)
	}
	return args.Get(0).(domain.Login), nil
}

func (m *LoginRepositoryMock) DeleteLogin(id int) error {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}
