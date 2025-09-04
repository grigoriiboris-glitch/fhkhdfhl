package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mymindmap/api/auth"
	"github.com/mymindmap/api/internal/handlers"
	"github.com/mymindmap/api/repository"
)

type Config struct {
	PostgresURL string
}

func loadConfig() (*Config, error) {
	_ = godotenv.Load()

	db := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")

	if db == "" || user == "" || pass == "" || host == "" {
		return nil, fmt.Errorf("missing database env vars")
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		user, pass, host, db,
	)

	return &Config{PostgresURL: connStr}, nil
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
	}
	defer dbpool.Close()

	// Репозитории
	postRepo := repository.NewPostRepository(dbpool)
	userRepo := repository.NewUserRepository(dbpool)
	mindMapRepo := repository.NewMindMapRepository(dbpool)

	// Сервисы
	authConfig := &auth.Config{
		EnableRateLimit: true,
	}
	authService, err := auth.NewAuthService(userRepo, authConfig)
	if err != nil {
		log.Fatalf("auth service error: %v", err)
	}

	// Handlers
	authHandler := handlers.NewAuthHandler(authService, userRepo, log.Default())
	postHandler := handlers.NewPostHandler(postRepo, authService, log.Default())
	mindMapHandler := handlers.NewMindMapHandler(mindMapRepo, authService, log.Default())

	// Router
	mux := http.NewServeMux()
	authHandler.RegisterRoutes(mux)
	postHandler.RegisterRoutes(mux)
	mindMapHandler.RegisterRoutes(mux)

	// Healthcheck
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	addr := ":8000"
	log.Printf("server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
