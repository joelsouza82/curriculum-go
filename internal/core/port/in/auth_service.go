package in

// AuthService é a porta de entrada (primária) para autenticação.
type AuthService interface {
	Authenticate(email, password string) (string, error)
}
