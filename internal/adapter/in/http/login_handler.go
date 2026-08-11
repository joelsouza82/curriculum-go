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

func (l *LoginHandler) GetLogins(ctx *gin.Context) {
	logins, err := l.loginService.GetLogins()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, logins)
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

	ctx.JSON(http.StatusOK, login)
}

func (l *LoginHandler) CreateLogin(ctx *gin.Context) {
	var login domain.Login
	err := ctx.BindJSON(&login)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertedLogin, err := l.loginService.CreateLogin(login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, insertedLogin)
}

func (l *LoginHandler) UpdateLogin(ctx *gin.Context) {
	idStr := ctx.Param("loginId")
	loginId, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de login inválido"})
		return
	}

	var login domain.Login
	if err := ctx.BindJSON(&login); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	login.ID = loginId

	updatedLogin, err := l.loginService.UpdateLogin(login)
	if err != nil {
		if errors.Is(err, domain.ErrLoginNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Registro de login não encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedLogin)
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
