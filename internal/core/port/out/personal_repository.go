package out

import "github.com/joelsouza82/curriculum-go/internal/core/domain"

// PersonalRepository é a porta de saída (secundária) para persistência de Personal.
type PersonalRepository interface {
	GetPersonalByID(id int) (domain.Personal, error)
	GetPersonals() ([]domain.Personal, error)
	CreatePersonal(personal domain.Personal) (int, error)
	UpdatePersonal(personal domain.Personal) (domain.Personal, error)
	DeletePersonal(id int) error
}
