package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	PostgresURL     string
	PostgresDB      string
	PostgresUser    string
	PostgresPass    string
	PostgresHost    string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:         getEnv("APP_PORT", "8000"),
		PostgresDB:   os.Getenv("POSTGRES_DB"),
		PostgresUser: os.Getenv("POSTGRES_USER"),
		PostgresPass: os.Getenv("POSTGRES_PASSWORD"),
		PostgresHost: os.Getenv("POSTGRES_HOST"),
	}

	cfg.PostgresURL = fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresHost, cfg.PostgresDB,
	)

	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
