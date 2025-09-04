package auth

import (
	"encoding/hex"
	"log/slog"
	"os"
	"strconv"
	"time"
)

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
	}

	// Session Key (required in production)
	if sessionKeyHex := os.Getenv("SESSION_KEY"); sessionKeyHex != "" {
		sessionKey, err := hex.DecodeString(sessionKeyHex)
		if err != nil {
			return nil, err
		}
		config.SessionKey = sessionKey
	}

	// Token expiration
	if tokenExpStr := os.Getenv("TOKEN_EXPIRATION_HOURS"); tokenExpStr != "" {
		hours, err := strconv.Atoi(tokenExpStr)
		if err != nil {
			return nil, err
		}
		config.TokenExpiration = time.Duration(hours) * time.Hour
	}

	// Refresh token expiration
	if refreshExpStr := os.Getenv("REFRESH_TOKEN_EXPIRATION_HOURS"); refreshExpStr != "" {
		hours, err := strconv.Atoi(refreshExpStr)
		if err != nil {
			return nil, err
		}
		config.RefreshTokenExp = time.Duration(hours) * time.Hour
	}

	// Bcrypt cost
	if bcryptCostStr := os.Getenv("BCRYPT_COST"); bcryptCostStr != "" {
		cost, err := strconv.Atoi(bcryptCostStr)
		if err != nil {
			return nil, err
		}
		config.BcryptCost = cost
	}

	// Rate limiting
	if rateLimitStr := os.Getenv("ENABLE_RATE_LIMIT"); rateLimitStr == "true" {
		config.EnableRateLimit = true
	}

	if maxAttemptsStr := os.Getenv("MAX_LOGIN_ATTEMPTS"); maxAttemptsStr != "" {
		attempts, err := strconv.Atoi(maxAttemptsStr)
		if err != nil {
			return nil, err
		}
		config.MaxLoginAttempts = attempts
	}

	if windowMinutesStr := os.Getenv("RATE_LIMIT_WINDOW_MINUTES"); windowMinutesStr != "" {
		minutes, err := strconv.Atoi(windowMinutesStr)
		if err != nil {
			return nil, err
		}
		config.RateLimitWindow = time.Duration(minutes) * time.Minute
	}

	if blockMinutesStr := os.Getenv("RATE_LIMIT_BLOCK_MINUTES"); blockMinutesStr != "" {
		minutes, err := strconv.Atoi(blockMinutesStr)
		if err != nil {
			return nil, err
		}
		config.RateLimitBlock = time.Duration(minutes) * time.Minute
	}

	return config, nil
}