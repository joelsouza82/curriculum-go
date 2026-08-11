package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type Config struct {
	ServerPort string
	Database   DatabaseConfig
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		ServerPort: getEnv("SERVER_PORT", "8000"),
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  getEnv("DB_SSLMODE", "require"),
		},
	}

	required := map[string]string{
		"DB_HOST":     cfg.Database.Host,
		"DB_PORT":     cfg.Database.Port,
		"DB_USER":     cfg.Database.User,
		"DB_PASSWORD": cfg.Database.Password,
		"DB_NAME":     cfg.Database.Name,
	}
	for name, value := range required {
		if value == "" {
			return Config{}, fmt.Errorf("variável de ambiente obrigatória não definida: %s", name)
		}
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
