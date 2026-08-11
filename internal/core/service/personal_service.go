package service

import (
	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	"github.com/joelsouza82/curriculum-go/internal/core/port/out"
)

type PersonalService struct {
	repository out.PersonalRepository
}

func NewPersonalService(repository out.PersonalRepository) *PersonalService {
	return &PersonalService{
		repository: repository,
	}
}

func (s *PersonalService) GetPersonalByID(id int) (domain.Personal, error) {
	return s.repository.GetPersonalByID(id)
}

func (s *PersonalService) GetPersonals() ([]domain.Personal, error) {
	return s.repository.GetPersonals()
}

func (s *PersonalService) CreatePersonal(personal domain.Personal) (domain.Personal, error) {
	personalId, err := s.repository.CreatePersonal(personal)
	if err != nil {
		return domain.Personal{}, err
	}
	personal.ID = personalId
	return personal, nil
}

func (s *PersonalService) UpdatePersonal(personal domain.Personal) (domain.Personal, error) {
	return s.repository.UpdatePersonal(personal)
}

func (s *PersonalService) DeletePersonal(id int) error {
	return s.repository.DeletePersonal(id)
}
