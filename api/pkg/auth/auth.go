// internal/auth/auth.go
package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrExpiredToken      = errors.New("token expired")
	ErrInvalidClaims     = errors.New("invalid token claims")
	ErrInsufficientRoles = errors.New("insufficient roles")
)

type Claims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type Config struct {
	JWTSecret         string
	AuthServiceURL    string
	TokenRefreshInterval time.Duration
}

type AuthService struct {
	config    Config
	jwtSecret []byte
	mu        sync.RWMutex
	userCache map[string]UserInfo // Кэш информации о пользователях
}

type UserInfo struct {
	UserID    string
	Role      string
	Permissions []string
	ExpiresAt time.Time
}

func NewAuthService(config Config) (*AuthService, error) {
	if config.JWTSecret == "" {
		return nil, errors.New("JWT secret is required")
	}
	
	if config.TokenRefreshInterval == 0 {
		config.TokenRefreshInterval = 5 * time.Minute
	}

	service := &AuthService{
		config:    config,
		jwtSecret: []byte(config.JWTSecret),
		userCache: make(map[string]UserInfo),
	}

	// Запускаем горутину для очистки кэша
	go service.cleanupCache()

	return service, nil
}

// ParseAndValidateToken парсит и валидирует JWT токен
func (a *AuthService) ParseAndValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GetUserInfo получает информацию о пользователе (из кэша или auth service)
func (a *AuthService) GetUserInfo(userID string) (UserInfo, error) {
	// Пытаемся получить из кэша
	if info, found := a.getFromCache(userID); found {
		return info, nil
	}

	// Если нет в кэше, запрашиваем у сервиса авторизации
	info, err := a.fetchUserInfoFromAuthService(userID)
	if err != nil {
		return UserInfo{}, err
	}

	// Сохраняем в кэш
	a.saveToCache(userID, info)

	return info, nil
}

// HasRole проверяет, есть ли у пользователя нужная роль
func (a *AuthService) HasRole(userID string, requiredRole string) (bool, error) {
	info, err := a.GetUserInfo(userID)
	if err != nil {
		return false, err
	}

	return info.Role == requiredRole, nil
}

// HasAnyRole проверяет, есть ли у пользователя любая из требуемых ролей
func (a *AuthService) HasAnyRole(userID string, requiredRoles []string) (bool, error) {
	info, err := a.GetUserInfo(userID)
	if err != nil {
		return false, err
	}

	for _, role := range requiredRoles {
		if info.Role == role {
			return true, nil
		}
	}

	return false, nil
}

// HasPermission проверяет, есть ли у пользователя определенное право
func (a *AuthService) HasPermission(userID string, permission string) (bool, error) {
	info, err := a.GetUserInfo(userID)
	if err != nil {
		return false, err
	}

	for _, perm := range info.Permissions {
		if perm == permission {
			return true, nil
		}
	}

	return false, nil
}

// fetchUserInfoFromAuthService запрашивает информацию о пользователе из сервиса авторизации
func (a *AuthService) fetchUserInfoFromAuthService(userID string) (UserInfo, error) {
	// В реальной реализации здесь должен быть HTTP запрос к сервису авторизации
	// Для примера используем заглушку
	
	// Пример реализации:
	/*
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", a.config.AuthServiceURL+"/user/"+userID, nil)
	resp, err := client.Do(req)
	if err != nil {
		return UserInfo{}, err
	}
	defer resp.Body.Close()
	
	// Парсим ответ
	*/
	
	// Заглушка: возвращаем тестовые данные в зависимости от userID
	switch userID {
	case "admin":
		return UserInfo{
			UserID:    "admin",
			Role:      "admin",
			Permissions: []string{"read", "write", "delete", "admin"},
			ExpiresAt: time.Now().Add(a.config.TokenRefreshInterval),
		}, nil
	case "user1":
		return UserInfo{
			UserID:    "user1",
			Role:      "user",
			Permissions: []string{"read", "write"},
			ExpiresAt: time.Now().Add(a.config.TokenRefreshInterval),
		}, nil
	default:
		return UserInfo{
			UserID:    userID,
			Role:      "guest",
			Permissions: []string{"read"},
			ExpiresAt: time.Now().Add(a.config.TokenRefreshInterval),
		}, nil
	}
}

// getFromCache получает информацию из кэша
func (a *AuthService) getFromCache(userID string) (UserInfo, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	info, exists := a.userCache[userID]
	if !exists || time.Now().After(info.ExpiresAt) {
		return UserInfo{}, false
	}

	return info, true
}

// saveToCache сохраняет информацию в кэш
func (a *AuthService) saveToCache(userID string, info UserInfo) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.userCache[userID] = info
}

// cleanupCache периодически очищает просроченный кэш
func (a *AuthService) cleanupCache() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		a.mu.Lock()
		now := time.Now()
		for userID, info := range a.userCache {
			if now.After(info.ExpiresAt) {
				delete(a.userCache, userID)
			}
		}
		a.mu.Unlock()
	}
}

// ExtractTokenFromHeader извлекает токен из заголовка Authorization
func ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header format must be: Bearer {token}")
	}

	return parts[1], nil
}