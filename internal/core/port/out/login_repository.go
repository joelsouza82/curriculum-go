package out

import "github.com/joelsouza82/curriculum-go/internal/core/domain"

// LoginRepository é a porta de saída (secundária) para persistência de Login.
type LoginRepository interface {
	GetLogins() ([]domain.Login, error)
	GetLoginByID(id int) (domain.Login, error)
	CreateLogin(login domain.Login) (int, error)
	UpdateLogin(login domain.Login) (domain.Login, error)
	DeleteLogin(id int) error
}
