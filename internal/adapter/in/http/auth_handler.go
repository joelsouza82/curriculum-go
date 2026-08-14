package http

import (
	"errors"
	"net/http"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	inport "github.com/joelsouza82/curriculum-go/internal/core/port/in"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService inport.AuthService
}

func NewAuthHandler(service inport.AuthService) AuthHandler {
	return AuthHandler{
		authService: service,
	}
}

type authRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *AuthHandler) Login(ctx *gin.Context) {
	var req authRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.authService.Authenticate(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "credenciais inválidas"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
