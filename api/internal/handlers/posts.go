package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/mymindmap/api/auth"
	"github.com/mymindmap/api/models"
	"github.com/mymindmap/api/repository"
)

type PostHandler struct {
	postRepo    *repository.PostRepository
	authService *auth.AuthService
	logger      *log.Logger
}

func NewPostHandler(postRepo *repository.PostRepository, authService *auth.AuthService, logger *log.Logger) *PostHandler {
	return &PostHandler{
		postRepo:    postRepo,
		authService: authService,
		logger:      logger,
	}
}

func (h *PostHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/posts", h.authService.AuthMiddleware(h.handlePosts))       // GET list, POST create
	mux.HandleFunc("/api/posts/", h.authService.AuthMiddleware(h.handleSinglePost)) // GET, PUT, DELETE by id
}

// --- Handlers ---

// handlePosts -> /api/posts
func (h *PostHandler) handlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetPosts(w, r)
	case http.MethodPost:
		h.CreatePost(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSinglePost -> /api/posts/{id}
func (h *PostHandler) handleSinglePost(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetPost(w, r, id)
	case http.MethodPut:
		h.UpdatePost(w, r, id)
	case http.MethodDelete:
		h.DeletePost(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetPosts - список постов
func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postRepo.GetAllPosts(r.Context())
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, posts)
}

// GetPost - один пост
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request, id int) {
	post, err := h.postRepo.GetPostByID(r.Context(), id)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if post == nil {
		h.respondError(w, http.StatusNotFound, "post not found")
		return
	}
	h.respondJSON(w, http.StatusOK, post)
}

// CreatePost - создание
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Title == "" || req.Content == "" {
		h.respondError(w, http.StatusBadRequest, "title and content required")
		return
	}

	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	post := &models.Post{Title: req.Title, Content: req.Content, UserID: user.UserID}
	if err := h.postRepo.CreatePost(r.Context(), post); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, post)
}

// UpdatePost - обновление
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request, id int) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	post, err := h.postRepo.GetPostByID(r.Context(), id)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if post == nil {
		h.respondError(w, http.StatusNotFound, "post not found")
		return
	}
	if post.UserID != user.UserID && user.Role != "admin" {
		h.respondError(w, http.StatusForbidden, "forbidden")
		return
	}

	post.Title, post.Content = req.Title, req.Content
	if err := h.postRepo.UpdatePost(r.Context(), post); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, post)
}

// DeletePost - удаление
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request, id int) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	post, err := h.postRepo.GetPostByID(r.Context(), id)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if post == nil {
		h.respondError(w, http.StatusNotFound, "post not found")
		return
	}
	if post.UserID != user.UserID && user.Role != "admin" {
		h.respondError(w, http.StatusForbidden, "forbidden")
		return
	}

	if err := h.postRepo.DeletePost(r.Context(), id); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]any{"success": true})
}

// --- Helpers ---

func (h *PostHandler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Printf("json encode error: %v", err)
	}
}

func (h *PostHandler) respondError(w http.ResponseWriter, status int, msg string) {
	h.respondJSON(w, status, map[string]string{"error": msg})
}
