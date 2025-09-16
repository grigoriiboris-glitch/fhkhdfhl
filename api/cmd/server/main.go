package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/go-chi/chi/v5"
	"github.com/mymindmap/api/internal/auth"
	"github.com/mymindmap/api/internal/http/routes"
)

type Config struct {
	PostgresURL string
	JWTSecret   string
}

func loadConfig() (*Config, error) {
	_ = godotenv.Load()

	db := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	jwtSecret := os.Getenv("JWT_SECRET")

	if db == "" || user == "" || pass == "" || host == "" {
		return nil, fmt.Errorf("missing database env vars")
	}

	if jwtSecret == "" {
		jwtSecret = "default-jwt-secret-change-in-production"
		//log.Println("WARNING: Using default JWT secret. Set JWT_SECRET env variable in production!")
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		user, pass, host, db,
	)

	return &Config{
		PostgresURL: connStr,
		JWTSecret:   jwtSecret,
	}, nil
}

func main() {
	// Загружаем конфиг
	conf, err := loadConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// Логирование в файл + stdout
	if err := os.MkdirAll("logs", 0755); err == nil {
		logFilePath := filepath.Join("logs", "server.log")
		if f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			log.SetOutput(f)
			defer f.Close()
		}
	}

	// Подключение к БД
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, conf.PostgresURL)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
		log.Println("db connection error: %v", err)
	}
	defer dbpool.Close()

	// Конфигурация аутентификации
	authConfig, err := auth.NewConfigFromEnv(slog.Default())
	if err != nil {
		log.Fatalf("auth config error: %v", err)
	}
	if conf.JWTSecret != "" {
		authConfig.JWTSecret = []byte(conf.JWTSecret)
	}

	// Роутер через DI
	// Получаем основной роутер с CRUD
	mainRouter, err := routes.NewRouter(ctx, dbpool, authConfig)
	if err != nil {
		log.Fatalf("router error: %v", err)
	}

	// Оборачиваем основной роутер в chi.Router для добавления ручных маршрутов
	r := chi.NewRouter()


	r.Mount("/api", mainRouter) // подключаем все CRUD маршруты


	addr := ":8000"
	log.Printf("server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
