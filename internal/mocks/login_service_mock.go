package mocks

import (
	"github.com/joelsouza82/curriculum-go/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type LoginServiceMock struct {
	mock.Mock
}

func (m *LoginServiceMock) GetLogins() ([]domain.Login, error) {
	args := m.Called()
	if args.Get(1) != nil {
		return []domain.Login{}, args.Error(1)
	}
	return args.Get(0).([]domain.Login), nil
}

func (m *LoginServiceMock) GetLoginByID(id int) (domain.Login, error) {
	args := m.Called(id)
	if args.Get(1) != nil {
		return domain.Login{}, args.Error(1)
	}
	return args.Get(0).(domain.Login), nil
}

func (m *LoginServiceMock) CreateLogin(login domain.Login) (domain.Login, error) {
	args := m.Called(login)
	if args.Get(1) != nil {
		return domain.Login{}, args.Error(1)
	}
	return args.Get(0).(domain.Login), nil
}

func (m *LoginServiceMock) UpdateLogin(login domain.Login) (domain.Login, error) {
	args := m.Called(login)
	if args.Get(1) != nil {
		return domain.Login{}, args.Error(1)
	}
	return args.Get(0).(domain.Login), nil
}

func (m *LoginServiceMock) DeleteLogin(id int) error {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}
