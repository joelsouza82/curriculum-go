package in

import "github.com/joelsouza82/curriculum-go/internal/core/domain"

// LoginService é a porta de entrada (primária) para os casos de uso de Login.
type LoginService interface {
	GetLogins() ([]domain.Login, error)
	GetLoginByID(id int) (domain.Login, error)
	CreateLogin(login domain.Login) (domain.Login, error)
	UpdateLogin(login domain.Login) (domain.Login, error)
	DeleteLogin(id int) error
}
