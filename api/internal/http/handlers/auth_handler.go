package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/mymindmap/api/internal/auth"
	"github.com/mymindmap/api/internal/http/middleware"
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
	mux.HandleFunc("/auth/refresh", h.RefreshToken)
	mux.HandleFunc("/auth/check", middleware.AuthMiddleware(h.authService, h.Check))
	mux.HandleFunc("/auth/user", middleware.AuthMiddleware(h.authService, h.GetCurrentUser))
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
	tokenPair, err := h.authService.LoginUser(r.Context(), req)
	if err != nil {
		h.respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	config := h.authService.GetConfig()
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenPair.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Установите true для HTTPS
		MaxAge:   int(config.TokenExpiration.Seconds()),
	})

	// Set refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Установите true для HTTPS
		MaxAge:   int(config.RefreshTokenExp.Seconds()),
	})
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
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
		// Clear access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	})

	// Clear refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	})

	h.respondJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get refresh token from cookie or header
	var refreshToken string
	
	// Try cookie first
	if cookie, err := r.Cookie("refresh_token"); err == nil && cookie.Value != "" {
		refreshToken = cookie.Value
	} else {
		// Try Authorization header
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			refreshToken = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if refreshToken == "" {
		h.respondError(w, http.StatusBadRequest, "refresh token is required")
		return
	}

	tokenPair, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		h.respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	config := h.authService.GetConfig()
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenPair.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Установите true для HTTPS
		MaxAge:   int(config.TokenExpiration.Seconds()),
	})

	// Set refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Установите true для HTTPS
		MaxAge:   int(config.RefreshTokenExp.Seconds()),
	})
	
	h.respondJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"access_token": tokenPair.AccessToken,
		"expires_at":   tokenPair.ExpiresAt,
	})
}

func (h *AuthHandler) Check(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("auth service error:")
	claims := middleware.GetUserFromContext(r.Context())
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
