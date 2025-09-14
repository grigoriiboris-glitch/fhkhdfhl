package routes

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mymindmap/api/internal/auth"
	"github.com/mymindmap/api/internal/http/handlers"
	"github.com/mymindmap/api/internal/services"
	"github.com/mymindmap/api/repository"
)

func NewRouter(ctx context.Context, dbpool *pgxpool.Pool, authConfig *auth.Config) (http.Handler, error) {
	r := chi.NewRouter()

	// ===== Repositories =====
	logRepo := repository.NewLogRepository(dbpool)

	// ===== Services =====
	logService := services.NewLogService(logRepo)

	// ===== Handlers =====
	logHandler := handlers.NewLogHandler(logService)

	// ===== RegisterRoutes =====
	r.Route("/logs", func(r chi.Router) {
		r.Get("/", logHandler.List)
		r.Post("/", logHandler.Create)
		r.Get("/{id}", logHandler.Get)
		r.Put("/{id}", logHandler.Update)
		r.Delete("/{id}", logHandler.Delete)
	})

	// ===== Healthcheck =====
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return r, nil
}
