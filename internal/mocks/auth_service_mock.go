package mocks

import (
	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) Authenticate(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}
