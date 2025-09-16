package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mymindmap/api/internal/auth"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

// AuthMiddleware проверяет JWT токен и добавляет пользователя в контекст
func AuthMiddleware(authService *auth.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		 if authService == nil {
            http.Error(w, "Authentication service unavailable", http.StatusInternalServerError)
            return
        }
		// Получаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Пробуем получить токен из cookie
			cookie, err := r.Cookie("auth_token")
			if err != nil || cookie.Value == "" {
				// Если нет токена, продолжаем без авторизации
				next.ServeHTTP(w, r)
				return
			}
			authHeader = "Bearer " + cookie.Value
		}

		// Проверяем формат заголовка
		if !strings.HasPrefix(authHeader, "Bearer ") {
			next.ServeHTTP(w, r)
			return
		}

		// Извлекаем токен
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Валидируем токен
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			// Если токен невалиден, продолжаем без авторизации
			next.ServeHTTP(w, r)
			return
		}

		// Добавляем пользователя в контекст
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// RequireAuth middleware требует авторизации для доступа
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(UserContextKey)
		if claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// RequirePermission middleware проверяет права доступа
func RequirePermission(authService *auth.AuthService, object, action string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(UserContextKey)
			if claims == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userClaims := claims.(*auth.Claims)
			
			// Проверяем права доступа по роли пользователя
			if !authService.CheckPermission(userClaims.Role, object, action) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) *auth.Claims {
	if user, ok := ctx.Value(UserContextKey).(*auth.Claims); ok {
		return user
	}
	return nil
}