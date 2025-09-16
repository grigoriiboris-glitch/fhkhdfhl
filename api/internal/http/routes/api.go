package routes

import (
	"context"
	"net/http"
	"log"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mymindmap/api/internal/auth"
	"github.com/mymindmap/api/internal/http/handlers"
	"github.com/mymindmap/api/internal/services/log_service"
	"github.com/mymindmap/api/repository"
	"github.com/mymindmap/api/internal/http/middleware"
)

func NewRouter(ctx context.Context, dbpool *pgxpool.Pool, authConfig *auth.Config) (http.Handler, error) {
	r := chi.NewRouter()

	// ===== Repositories =====
	logRepo := repository.NewLogRepository(dbpool)
	userRepo := repository.NewUserRepository(dbpool)

	// ===== Services =====
	logService := log_service.NewLogService(logRepo)
	authService, err := auth.NewAuthService(userRepo, authConfig)

	if err != nil  { 
		log.Fatalf("auth service error: %v", err)
	}

	// ===== Handlers =====
	logHandler := handlers.NewLogHandler(logService)
	auth_handler := handlers.NewAuthHandler(authService,userRepo,log.Default())

	r.Route("/auth", func(r chi.Router) {
    r.Post("/login", auth_handler.Login)
    r.Post("/register", auth_handler.Register)
    r.Post("/logout", auth_handler.Logout)
    r.Post("/refresh", auth_handler.RefreshToken)
    r.Get("/check", middleware.AuthMiddleware(authService, auth_handler.Check))
    r.Get("/user", middleware.AuthMiddleware(authService, auth_handler.GetCurrentUser))
    })


	// ===== Manual routes at root level =====
	r.Get("/manual/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from manual route!"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"OK"}`))
	})

		r.Route("/logs", func(r chi.Router) {
			r.Get("/", logHandler.List)
			r.Post("/", logHandler.Create)
			r.Get("/{id}", logHandler.Get)
			r.Put("/{id}", logHandler.Update)
			r.Delete("/{id}", logHandler.Delete)
		})

		// Healthcheck under /api
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"api":"healthy"}`))
		})

	return r, nil
}