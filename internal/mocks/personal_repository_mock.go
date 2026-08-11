package mocks

import (
	"github.com/joelsouza82/curriculum-go/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type PersonalRepositoryMock struct {
	mock.Mock
}

func (m *PersonalRepositoryMock) GetPersonalByID(id int) (domain.Personal, error) {
	args := m.Called(id)
	if args.Get(1) != nil {
		return domain.Personal{}, args.Error(1)
	}
	return args.Get(0).(domain.Personal), nil
}

func (m *PersonalRepositoryMock) GetPersonals() ([]domain.Personal, error) {
	args := m.Called()
	if args.Get(1) != nil {
		return []domain.Personal{}, args.Error(1)
	}
	return args.Get(0).([]domain.Personal), nil
}

func (m *PersonalRepositoryMock) CreatePersonal(personal domain.Personal) (int, error) {
	args := m.Called(personal)
	if args.Get(1) != nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int), nil
}

func (m *PersonalRepositoryMock) UpdatePersonal(personal domain.Personal) (domain.Personal, error) {
	args := m.Called(personal)
	if args.Get(1) != nil {
		return domain.Personal{}, args.Error(1)
	}
	return args.Get(0).(domain.Personal), nil
}

func (m *PersonalRepositoryMock) DeletePersonal(id int) error {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}
