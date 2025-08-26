package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    AppPort    string
    AppEnv     string
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
    JWTSecret  string
    JWTExpire  string
}

func LoadConfig() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    return &Config{
        AppPort:    getEnv("APP_PORT", "8080"),
        AppEnv:     getEnv("APP_ENV", "development"),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "ecolokal"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
        JWTSecret:  getEnv("JWT_SECRET", "secret"),
        JWTExpire:  getEnv("JWT_EXPIRE_HOURS", "24"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}