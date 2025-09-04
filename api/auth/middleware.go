package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

// AuthMiddleware проверяет JWT токен и добавляет пользователя в контекст
func (s *AuthService) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		claims, err := s.ValidateToken(tokenString)
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
func (s *AuthService) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
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
func (s *AuthService) RequirePermission(object, action string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(UserContextKey)
			if claims == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userClaims := claims.(*Claims)
			
			// Проверяем права доступа по роли пользователя
			if !s.CheckPermission(userClaims.Role, object, action) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) *Claims {
	if user, ok := ctx.Value(UserContextKey).(*Claims); ok {
		return user
	}
	return nil
}

// SetAuthCookie устанавливает cookie с токеном авторизации
func (s *AuthService) SetAuthCookie(w http.ResponseWriter, tokenPair *TokenPair) {
	// Set access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenPair.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Установите true для HTTPS
		MaxAge:   int(s.config.TokenExpiration.Seconds()),
	})

	// Set refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Установите true для HTTPS
		MaxAge:   int(s.config.RefreshTokenExp.Seconds()),
	})
}

// ClearAuthCookie удаляет cookie авторизации
func (s *AuthService) ClearAuthCookie(w http.ResponseWriter) {
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
} 