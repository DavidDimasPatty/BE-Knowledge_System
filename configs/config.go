package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppEnv  string
	DBHost  string
	DBUser  string
	DBPass  string
	DBName  string
}

func LoadConfig() *Config {
	// Load .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, using system environment")
	}

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),
		DBHost:  getEnv("DB_HOST", "localhost"),
		DBUser:  getEnv("DB_USER", "root"),
		DBPass:  getEnv("DB_PASS", ""),
		DBName:  getEnv("DB_NAME", "knowledge_system"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
