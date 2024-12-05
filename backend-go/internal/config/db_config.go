package config

import (
	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSchema   string
}

func LoadDBConfig() (*DBConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &DBConfig{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "topics"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBSchema:   getEnv("DB_SCHEMA", "public"),
	}, nil
}
