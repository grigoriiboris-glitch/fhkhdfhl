package auth

import (
	"context"
	"log/slog"
	"os"

	"github.com/mymindmap/api/models"
	"github.com/mymindmap/api/repository"
)

// ExampleUsage demonstrates how to use the improved auth service
func ExampleUsage() {
	// Create logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create config from environment variables
	config, err := NewConfigFromEnv(logger)
	if err != nil {
		logger.Error("Failed to create config", "error", err)
		return
	}

	// Enable rate limiting for production
	config.EnableRateLimit = true

	// Create user repository (implementation depends on your database)
	userRepo := &repository.UserRepository{} // Initialize with your database connection

	// Create auth service
	authService, err := NewAuthService(userRepo, config)
	if err != nil {
		logger.Error("Failed to create auth service", "error", err)
		return
	}

	ctx := context.Background()

	// Register a new user
	registerReq := &models.RegisterRequest{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "SecureP@ssw0rd123!",
	}

	user, err := authService.RegisterUser(ctx, registerReq)
	if err != nil {
		logger.Error("Failed to register user", "error", err)
		return
	}

	logger.Info("User registered successfully", "user_id", user.ID, "email", user.Email)

	// Login user
	loginReq := &models.LoginRequest{
		Email:    "john.doe@example.com",
		Password: "SecureP@ssw0rd123!",
	}

	tokenPair, err := authService.LoginUser(ctx, loginReq)
	if err != nil {
		logger.Error("Failed to login user", "error", err)
		return
	}

	logger.Info("User logged in successfully", "expires_at", tokenPair.ExpiresAt)

	// Validate token
	claims, err := authService.ValidateToken(tokenPair.AccessToken)
	if err != nil {
		logger.Error("Failed to validate token", "error", err)
		return
	}

	logger.Info("Token validated", "user_id", claims.UserID, "email", claims.Email, "role", claims.Role)

	// Check permissions
	canRead := authService.CheckPermissionForUser(claims.Email, ObjectPost, ActionRead)
	canDelete := authService.CheckPermissionForUser(claims.Email, ObjectPost, ActionDelete)

	logger.Info("Permission check", "can_read_posts", canRead, "can_delete_posts", canDelete)

	// Refresh token
	newTokenPair, err := authService.RefreshToken(tokenPair.RefreshToken)
	if err != nil {
		logger.Error("Failed to refresh token", "error", err)
		return
	}

	logger.Info("Token refreshed successfully", "expires_at", newTokenPair.ExpiresAt)
}

// Environment variables example:
// JWT_SECRET=your_hex_encoded_secret_key_here
// SESSION_KEY=your_hex_encoded_session_key_here
// TOKEN_EXPIRATION_HOURS=24
// REFRESH_TOKEN_EXPIRATION_HOURS=168
// BCRYPT_COST=12
// ENABLE_RATE_LIMIT=true
// MAX_LOGIN_ATTEMPTS=5
// RATE_LIMIT_WINDOW_MINUTES=15
// RATE_LIMIT_BLOCK_MINUTES=15