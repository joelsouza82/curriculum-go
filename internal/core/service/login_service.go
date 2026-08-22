package service

import (
	"errors"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	"github.com/joelsouza82/curriculum-go/internal/core/port/out"

	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	repository out.LoginRepository
}

func NewLoginService(repository out.LoginRepository) *LoginService {
	return &LoginService{
		repository: repository,
	}
}

func (s *LoginService) GetLogins() ([]domain.Login, error) {
	return s.repository.GetLogins()
}

func (s *LoginService) GetLoginByID(id int) (domain.Login, error) {
	return s.repository.GetLoginByID(id)
}

func (s *LoginService) CreateLogin(login domain.Login) (domain.Login, error) {
	_, err := s.repository.GetLoginByEmail(login.Email)
	if err == nil {
		return domain.Login{}, domain.ErrEmailAlreadyExists
	}
	if !errors.Is(err, domain.ErrLoginNotFound) {
		return domain.Login{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.Login{}, err
	}
	login.Password = string(hashedPassword)

	loginId, err := s.repository.CreateLogin(login)
	if err != nil {
		return domain.Login{}, err
	}

	login.ID = loginId
	return login, nil
}

func (s *LoginService) UpdateLogin(login domain.Login) (domain.Login, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.Login{}, err
	}
	login.Password = string(hashedPassword)

	return s.repository.UpdateLogin(login)
}

func (s *LoginService) DeleteLogin(id int) error {
	return s.repository.DeleteLogin(id)
}
