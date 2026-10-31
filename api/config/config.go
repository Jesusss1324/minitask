package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port              string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	CORSAllowedOrigin string
}

func Load() Config {
	return Config{
		Port:              getEnv("APP_PORT", "8080"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "3306"),
		DBUser:            getEnv("DB_USER", "root"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            getEnv("DB_NAME", "minitask"),
		CORSAllowedOrigin: getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC&charset=utf8mb4",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}
