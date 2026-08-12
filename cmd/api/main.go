package main

import (
	"log"

	httpadapter "github.com/joelsouza82/curriculum-go/internal/adapter/in/http"
	"github.com/joelsouza82/curriculum-go/internal/adapter/out/postgres"
	"github.com/joelsouza82/curriculum-go/internal/config"
	"github.com/joelsouza82/curriculum-go/internal/core/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("erro ao carregar configuração: %v", err)
	}

	dbConnection, err := postgres.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("erro de conexão com o banco de dados: %v", err)
	}
	defer dbConnection.Close()

	personalRepository := postgres.NewPersonalRepository(dbConnection)
	personalService := service.NewPersonalService(personalRepository)
	personalHandler := httpadapter.NewPersonalHandler(personalService)

	loginRepository := postgres.NewLoginRepository(dbConnection)
	loginService := service.NewLoginService(loginRepository)
	loginHandler := httpadapter.NewLoginHandler(loginService)

	authService := service.NewAuthService(loginRepository, cfg.JWTSecret)
	authHandler := httpadapter.NewAuthHandler(authService)

	router := httpadapter.NewRouter(personalHandler, loginHandler, authHandler, cfg.JWTSecret)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
