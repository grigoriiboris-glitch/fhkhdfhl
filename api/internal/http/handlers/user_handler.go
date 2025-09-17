package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/mymindmap/api/internal/auth"
    "github.com/mymindmap/api/internal/http/middleware"
    "github.com/mymindmap/api/repository"
)

type UserHandler struct {
    userRepo    *repository.UserRepository
    authService *auth.AuthService
}

func NewUserHandler(userRepo *repository.UserRepository, authService *auth.AuthService) *UserHandler {
    return &UserHandler{userRepo: userRepo, authService: authService}
}

// AuthCheck returns 200 if authorized, 401 otherwise
func (h *UserHandler) AuthCheck(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    claims := middleware.GetUserFromContext(r.Context())
    w.Header().Set("Content-Type", "application/json")
    if claims == nil {
        w.WriteHeader(http.StatusUnauthorized)
        _ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
        return
    }
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

// GetCurrentUser returns current user as JSON
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    claims := middleware.GetUserFromContext(r.Context())
    if claims == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        _ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
        return
    }
    user, err := h.userRepo.GetUserByID(r.Context(), claims.UserID)
    if err != nil || user == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        _ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to get use handr"})
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(user)
}

