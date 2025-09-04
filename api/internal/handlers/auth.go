package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/mymindmap/api/auth"
	"github.com/mymindmap/api/models"
	"github.com/mymindmap/api/repository"
)

type AuthHandler struct {
	authService *auth.AuthService
	userRepo    *repository.UserRepository
	logger      *log.Logger
}

func NewAuthHandler(authService *auth.AuthService, userRepo *repository.UserRepository, logger *log.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, userRepo: userRepo, logger: logger}
}

// Регистрируем маршруты
func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/login", h.Login)
	mux.HandleFunc("/auth/register", h.Register)
	mux.HandleFunc("/auth/logout", h.Logout)
	mux.HandleFunc("/auth/check", h.authService.AuthMiddleware(h.Check))
	mux.HandleFunc("/auth/user", h.authService.AuthMiddleware(h.GetCurrentUser))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// поддержка JSON
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			h.respondError(w, http.StatusBadRequest, "invalid json")
			return
		}
	}

	if creds.Email == "" || creds.Password == "" {
		h.respondError(w, http.StatusBadRequest, "email и пароль обязательны")
		return
	}

	req := &models.LoginRequest{Email: creds.Email, Password: creds.Password}
	token, err := h.authService.LoginUser(r.Context(), req)
	if err != nil {
		h.respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	h.authService.SetAuthCookie(w, token)
	h.respondJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		h.respondError(w, http.StatusBadRequest, "все поля обязательны")
		return
	}

	if len(req.Password) < 6 {
		h.respondError(w, http.StatusBadRequest, "пароль должен содержать минимум 6 символов")
		return
	}

	user, err := h.authService.RegisterUser(r.Context(), &req)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.authService.ClearAuthCookie(w)
	h.respondJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *AuthHandler) Check(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.userRepo.GetUserByID(r.Context(), claims.UserID)
	if err != nil || user == nil {
		h.respondError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	h.respondJSON(w, http.StatusOK, user)
}

// --- helpers ---

func (h *AuthHandler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Printf("json encode error: %v", err)
	}
}

func (h *AuthHandler) respondError(w http.ResponseWriter, status int, msg string) {
	h.respondJSON(w, status, map[string]string{"error": msg})
}
