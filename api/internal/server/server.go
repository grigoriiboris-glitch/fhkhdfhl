package server

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymindmap/api/auth"
	"github.com/mymindmap/api/handlers"
	"github.com/mymindmap/api/internal/config"
	"github.com/mymindmap/api/repository"
)

type Server struct {
	http.Server
}

func New(cfg *config.Config, dbpool *pgxpool.Pool, log *log.Logger) *Server {
	postRepo := repository.NewPostRepository(dbpool)
	userRepo := repository.NewUserRepository(dbpool)
	mindMapRepo := repository.NewMindMapRepository(dbpool)

	authService, err := auth.NewAuthService(userRepo)
	if err != nil {
		log.Fatal("unable to init auth service:", err)
	}

	mux := http.NewServeMux()

	// auth routes
	authHandler := handlers.NewAuthHandler(authService, userRepo)
	authHandler.RegisterRoutes(mux)

	// post routes
	postHandler := handlers.NewPostHandler(postRepo, authService)
	postHandler.RegisterRoutes(mux)

	// mindmap routes
	mindMapHandler := handlers.NewMindMapHandler(mindMapRepo, authService)
	mindMapHandler.RegisterRoutes(mux)

	return &Server{
		Server: http.Server{
			Addr:    ":" + cfg.Port,
			Handler: loggingMiddleware(mux, log),
		},
	}
}
