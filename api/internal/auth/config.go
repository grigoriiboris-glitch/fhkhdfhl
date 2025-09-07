package auth

import (
	"encoding/hex"
	"log/slog"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Config конфигурация для сервиса аутентификации
type Config struct {
	JWTSecret        []byte        // Секрет для подписи JWT токенов
	SessionKey       []byte        // Ключ для сессий (если используются)
	TokenExpiration  time.Duration // Время жизни access токена
	RefreshTokenExp  time.Duration // Время жизни refresh токена
	BcryptCost       int           // Сложность bcrypt хеширования
	Logger           *slog.Logger  // Логгер для записи событий
	EnableRateLimit  bool          // Включить лимитирование запросов
	MaxLoginAttempts int           // Максимум попыток входа
	RateLimitWindow  time.Duration // Окно времени для лимита
	RateLimitBlock   time.Duration // Время блокировки после превышения лимита
}

// NewConfig создает новую конфигурацию с настройками по умолчанию
func NewConfig(logger *slog.Logger) *Config {
	return &Config{
		JWTSecret:        []byte("default-jwt-secret-change-in-production"),
		SessionKey:       []byte("default-session-key-change-in-production"),
		TokenExpiration:  15 * time.Minute,
		RefreshTokenExp:  7 * 24 * time.Hour,
		BcryptCost:       bcrypt.DefaultCost,
		Logger:           logger,
		EnableRateLimit:  true,
		MaxLoginAttempts: 5,
		RateLimitWindow:  15 * time.Minute,
		RateLimitBlock:   15 * time.Minute,
	}
}

// NewConfigFromEnv creates a new Config from environment variables
func NewConfigFromEnv(logger *slog.Logger) (*Config, error) {
	config := &Config{
		Logger: logger,
	}

	// JWT Secret (required in production)
	if jwtSecretHex := os.Getenv("JWT_SECRET"); jwtSecretHex != "" {
		jwtSecret, err := hex.DecodeString(jwtSecretHex)
		if err != nil {
			return nil, err
		}
		config.JWTSecret = jwtSecret
	} else {
		config.JWTSecret = []byte("default-jwt-secret-change-in-production")
	}

	// Session Key (required in production)
	if sessionKeyHex := os.Getenv("SESSION_KEY"); sessionKeyHex != "" {
		sessionKey, err := hex.DecodeString(sessionKeyHex)
		if err != nil {
			return nil, err
		}
		config.SessionKey = sessionKey
	} else {
		config.SessionKey = []byte("default-session-key-change-in-production")
	}

	// Token expiration
	if tokenExpStr := os.Getenv("TOKEN_EXPIRATION_HOURS"); tokenExpStr != "" {
		hours, err := strconv.Atoi(tokenExpStr)
		if err != nil {
			return nil, err
		}
		config.TokenExpiration = time.Duration(hours) * time.Hour
	} else {
		config.TokenExpiration = 15 * time.Minute
	}

	// Refresh token expiration
	if refreshExpStr := os.Getenv("REFRESH_TOKEN_EXPIRATION_HOURS"); refreshExpStr != "" {
		hours, err := strconv.Atoi(refreshExpStr)
		if err != nil {
			return nil, err
		}
		config.RefreshTokenExp = time.Duration(hours) * time.Hour
	} else {
		config.RefreshTokenExp = 7 * 24 * time.Hour
	}

	// Bcrypt cost
	if bcryptCostStr := os.Getenv("BCRYPT_COST"); bcryptCostStr != "" {
		cost, err := strconv.Atoi(bcryptCostStr)
		if err != nil {
			return nil, err
		}
		config.BcryptCost = cost
	} else {
		config.BcryptCost = bcrypt.DefaultCost
	}

	// Rate limiting
	if rateLimitStr := os.Getenv("ENABLE_RATE_LIMIT"); rateLimitStr == "true" {
		config.EnableRateLimit = true
	} else {
		config.EnableRateLimit = false
	}

	if maxAttemptsStr := os.Getenv("MAX_LOGIN_ATTEMPTS"); maxAttemptsStr != "" {
		attempts, err := strconv.Atoi(maxAttemptsStr)
		if err != nil {
			return nil, err
		}
		config.MaxLoginAttempts = attempts
	} else {
		config.MaxLoginAttempts = 5
	}

	if windowMinutesStr := os.Getenv("RATE_LIMIT_WINDOW_MINUTES"); windowMinutesStr != "" {
		minutes, err := strconv.Atoi(windowMinutesStr)
		if err != nil {
			return nil, err
		}
		config.RateLimitWindow = time.Duration(minutes) * time.Minute
	} else {
		config.RateLimitWindow = 1 * time.Minute
	}

	if blockMinutesStr := os.Getenv("RATE_LIMIT_BLOCK_MINUTES"); blockMinutesStr != "" {
		minutes, err := strconv.Atoi(blockMinutesStr)
		if err != nil {
			return nil, err
		}
		config.RateLimitBlock = time.Duration(minutes) * time.Minute
	} else {
		config.RateLimitBlock = 15 * time.Minute
	}

	return config, nil
}
