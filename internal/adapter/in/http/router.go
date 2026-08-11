package http

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(personalHandler PersonalHandler, loginHandler LoginHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/personals", personalHandler.GetPersonals)
	router.GET("/personal/:personalId", personalHandler.GetPersonalByID)
	router.POST("/personal", personalHandler.CreatePersonal)
	router.PUT("/personal/:personalId", personalHandler.UpdatePersonal)
	router.DELETE("/personal/:personalId", personalHandler.DeletePersonal)

	router.GET("/logins", loginHandler.GetLogins)
	router.GET("/login/:loginId", loginHandler.GetLoginByID)
	router.POST("/login", loginHandler.CreateLogin)
	router.PUT("/login/:loginId", loginHandler.UpdateLogin)
	router.DELETE("/login/:loginId", loginHandler.DeleteLogin)

	return router
}
