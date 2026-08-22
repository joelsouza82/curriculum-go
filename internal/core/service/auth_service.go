package service

import (
	"errors"
	"time"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	"github.com/joelsouza82/curriculum-go/internal/core/port/out"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const tokenExpiration = 1 * time.Hour

type AuthService struct {
	repository out.LoginRepository
	jwtSecret  []byte
}

func NewAuthService(repository out.LoginRepository, jwtSecret string) *AuthService {
	return &AuthService{
		repository: repository,
		jwtSecret:  []byte(jwtSecret),
	}
}

func (s *AuthService) Authenticate(email, password string) (string, error) {
	login, err := s.repository.GetLoginByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrLoginNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub":   login.ID,
		"email": login.Email,
		"exp":   time.Now().Add(tokenExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
