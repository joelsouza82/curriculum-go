package in

import "github.com/joelsouza82/curriculum-go/internal/core/domain"

// PersonalService é a porta de entrada (primária) para os casos de uso de Personal.
type PersonalService interface {
	GetPersonalByID(id int) (domain.Personal, error)
	GetPersonals() ([]domain.Personal, error)
	CreatePersonal(personal domain.Personal) (domain.Personal, error)
	UpdatePersonal(personal domain.Personal) (domain.Personal, error)
	DeletePersonal(id int) error
}
