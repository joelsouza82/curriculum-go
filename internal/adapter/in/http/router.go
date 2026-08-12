package http

import (
	"github.com/joelsouza82/curriculum-go/internal/adapter/in/http/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(personalHandler PersonalHandler, loginHandler LoginHandler, authHandler AuthHandler, jwtSecret string) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/auth/login", authHandler.Login)
	router.POST("/login", loginHandler.CreateLogin)

	protected := router.Group("/")
	protected.Use(middleware.RequireAuth(jwtSecret))
	{
		protected.GET("/personals", personalHandler.GetPersonals)
		protected.GET("/personal/:personalId", personalHandler.GetPersonalByID)
		protected.POST("/personal", personalHandler.CreatePersonal)
		protected.PUT("/personal/:personalId", personalHandler.UpdatePersonal)
		protected.DELETE("/personal/:personalId", personalHandler.DeletePersonal)

		protected.GET("/logins", loginHandler.GetLogins)
		protected.GET("/login/:loginId", loginHandler.GetLoginByID)
		protected.PUT("/login/:loginId", loginHandler.UpdateLogin)
		protected.DELETE("/login/:loginId", loginHandler.DeleteLogin)
	}

	return router
}
