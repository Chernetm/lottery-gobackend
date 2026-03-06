package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL              string
	JWTSecret          string
	Port               string
	ChapaSecret        string
	ChapaWebhookSecret string
	BaseURL            string
	FrontendURL        string
}

var AppConfig *Config

func LoadConfig() {
	// Try loading from current directory
	err := godotenv.Load()
	if err != nil {
		// Try loading from parent directory (useful if running from cmd/server)
		err = godotenv.Load("../../.env")
		if err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	AppConfig = &Config{
		DBURL:              getEnv("DATABASE_URL", "root:password@tcp(127.0.0.1:3306)/lottery?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTSecret:          getEnv("JWT_SECRET", "your_secret_key"),
		Port:               getEnv("PORT", "8081"),
		ChapaSecret:        getEnv("CHAPA_SECRET_KEY", ""),
		ChapaWebhookSecret: getEnv("CHAPA_WEBHOOK_SECRET", ""),
		BaseURL:            getEnv("BASE_URL", "http://localhost:5000"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
