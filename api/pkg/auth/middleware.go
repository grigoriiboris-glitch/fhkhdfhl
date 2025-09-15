// internal/auth/middleware.go
package auth

import (
	"context"
	"net/http"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

// Middleware создает middleware для проверки аутентификации
func (a *AuthService) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := ExtractTokenFromHeader(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := a.ParseAndValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Добавляем информацию о пользователе в контекст
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRoleMiddleware создает middleware для проверки ролей
func (a *AuthService) RequireRoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(UserIDKey).(string)
			if !ok {
				http.Error(w, "User not authenticated", http.StatusUnauthorized)
				return
			}

			hasRole, err := a.HasRole(userID, requiredRole)
			if err != nil {
				http.Error(w, "Error checking role: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if !hasRole {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRoleMiddleware создает middleware для проверки любой из ролей
func (a *AuthService) RequireAnyRoleMiddleware(requiredRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(UserIDKey).(string)
			if !ok {
				http.Error(w, "User not authenticated", http.StatusUnauthorized)
				return
			}

			hasAnyRole, err := a.HasAnyRole(userID, requiredRoles)
			if err != nil {
				http.Error(w, "Error checking roles: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if !hasAnyRole {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext извлекает информацию о пользователе из контекста
func GetUserFromContext(ctx context.Context) (userID, role string) {
	userID, _ = ctx.Value(UserIDKey).(string)
	role, _ = ctx.Value(RoleKey).(string)
	return
}