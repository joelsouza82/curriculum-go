package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	inport "github.com/joelsouza82/curriculum-go/internal/core/port/in"

	"github.com/gin-gonic/gin"
)

type PersonalHandler struct {
	personalService inport.PersonalService
}

func NewPersonalHandler(service inport.PersonalService) PersonalHandler {
	return PersonalHandler{
		personalService: service,
	}
}

func (p *PersonalHandler) GetPersonalByID(ctx *gin.Context) {
	idStr := ctx.Param("personalId")
	personalId, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de personal inválido"})
		return
	}

	personal, err := p.personalService.GetPersonalByID(personalId)
	if err != nil {
		if errors.Is(err, domain.ErrPersonalNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Registro de personal não encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, personal)
}

func (p *PersonalHandler) GetPersonals(ctx *gin.Context) {
	personals, err := p.personalService.GetPersonals()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, personals)
}

func (p *PersonalHandler) CreatePersonal(ctx *gin.Context) {
	var personal domain.Personal
	err := ctx.BindJSON(&personal)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertedPersonal, err := p.personalService.CreatePersonal(personal)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, insertedPersonal)
}

func (p *PersonalHandler) UpdatePersonal(ctx *gin.Context) {
	idStr := ctx.Param("personalId")
	personalId, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de personal inválido"})
		return
	}

	var personal domain.Personal
	if err := ctx.BindJSON(&personal); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	personal.ID = personalId

	updatedPersonal, err := p.personalService.UpdatePersonal(personal)
	if err != nil {
		if errors.Is(err, domain.ErrPersonalNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Registro de personal não encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedPersonal)
}

func (p *PersonalHandler) DeletePersonal(ctx *gin.Context) {
	idStr := ctx.Param("personalId")
	personalId, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de personal inválido"})
		return
	}

	if err := p.personalService.DeletePersonal(personalId); err != nil {
		if errors.Is(err, domain.ErrPersonalNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Registro de personal não encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
