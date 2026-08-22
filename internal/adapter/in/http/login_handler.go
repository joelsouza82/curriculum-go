package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	inport "github.com/joelsouza82/curriculum-go/internal/core/port/in"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	loginService inport.LoginService
}

func NewLoginHandler(service inport.LoginService) LoginHandler {
	return LoginHandler{
		loginService: service,
	}
}

// loginRequest é o payload aceito para criar/atualizar um login.
type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// loginResponse é o payload devolvido ao cliente; nunca inclui a senha (nem o hash).
type loginResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func toLoginResponse(login domain.Login) loginResponse {
	return loginResponse{ID: login.ID, Email: login.Email}
}

func (l *LoginHandler) GetLogins(ctx *gin.Context) {
	logins, err := l.loginService.GetLogins()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]loginResponse, 0, len(logins))
	for _, login := range logins {
		response = append(response, toLoginResponse(login))
	}
	ctx.JSON(http.StatusOK, response)
}

func (l *LoginHandler) GetLoginByID(ctx *gin.Context) {
	idStr := ctx.Param("loginId")
	loginId, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de login inválido"})
		return
	}

	login, err := l.loginService.GetLoginByID(loginId)
	if err != nil {
		if errors.Is(err, domain.ErrLoginNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Registro de login não encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, toLoginResponse(login))
}

func (l *LoginHandler) CreateLogin(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertedLogin, err := l.loginService.CreateLogin(domain.Login{Email: req.Email, Password: req.Password})
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "e-mail já cadastrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, toLoginResponse(insertedLogin))
}

func (l *LoginHandler) UpdateLogin(ctx *gin.Context) {
	idStr := ctx.Param("loginId")
	loginId, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de login inválido"})
		return
	}

	var req loginRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedLogin, err := l.loginService.UpdateLogin(domain.Login{ID: loginId, Email: req.Email, Password: req.Password})
	if err != nil {
		if errors.Is(err, domain.ErrLoginNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Registro de login não encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, toLoginResponse(updatedLogin))
}

func (l *LoginHandler) DeleteLogin(ctx *gin.Context) {
	idStr := ctx.Param("loginId")
	loginId, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de login inválido"})
		return
	}

	if err := l.loginService.DeleteLogin(loginId); err != nil {
		if errors.Is(err, domain.ErrLoginNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Registro de login não encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
