package auth

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ConfigTestSuite defines the test suite for Config
type ConfigTestSuite struct {
	suite.Suite
	originalEnv map[string]string
	logger      *slog.Logger
}

// SetupTest runs before each test
func (suite *ConfigTestSuite) SetupTest() {
	// Save original environment variables
	suite.originalEnv = make(map[string]string)
	envVars := []string{
		"JWT_SECRET",
		"SESSION_KEY",
		"TOKEN_EXPIRATION_HOURS",
		"REFRESH_TOKEN_EXPIRATION_HOURS",
		"BCRYPT_COST",
		"ENABLE_RATE_LIMIT",
		"MAX_LOGIN_ATTEMPTS",
		"RATE_LIMIT_WINDOW_MINUTES",
		"RATE_LIMIT_BLOCK_MINUTES",
	}
	
	for _, envVar := range envVars {
		suite.originalEnv[envVar] = os.Getenv(envVar)
		os.Unsetenv(envVar)
	}

	suite.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
}

// TearDownTest runs after each test
func (suite *ConfigTestSuite) TearDownTest() {
	// Restore original environment variables
	for envVar, value := range suite.originalEnv {
		if value == "" {
			os.Unsetenv(envVar)
		} else {
			os.Setenv(envVar, value)
		}
	}
}

// Test NewConfigFromEnv with no environment variables
func (suite *ConfigTestSuite) TestNewConfigFromEnv_NoEnvVars() {
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Equal(suite.T(), suite.logger, config.Logger)
	assert.Empty(suite.T(), config.JWTSecret)
	assert.Empty(suite.T(), config.SessionKey)
	assert.Zero(suite.T(), config.TokenExpiration)
	assert.Zero(suite.T(), config.RefreshTokenExp)
	assert.Zero(suite.T(), config.BcryptCost)
	assert.False(suite.T(), config.EnableRateLimit)
	assert.Zero(suite.T(), config.MaxLoginAttempts)
	assert.Zero(suite.T(), config.RateLimitWindow)
	assert.Zero(suite.T(), config.RateLimitBlock)
}

// Test NewConfigFromEnv with JWT_SECRET
func (suite *ConfigTestSuite) TestNewConfigFromEnv_WithJWTSecret() {
	// Set hex-encoded secret (32 bytes = 64 hex chars)
	jwtSecretHex := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	os.Setenv("JWT_SECRET", jwtSecretHex)
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Len(suite.T(), config.JWTSecret, 32)
	
	// Verify the decoded bytes match
	expectedBytes := []byte{
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
	}
	assert.Equal(suite.T(), expectedBytes, config.JWTSecret)
}

// Test NewConfigFromEnv with invalid JWT_SECRET
func (suite *ConfigTestSuite) TestNewConfigFromEnv_InvalidJWTSecret() {
	os.Setenv("JWT_SECRET", "invalid-hex-string")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), config)
	assert.Contains(suite.T(), err.Error(), "encoding/hex")
}

// Test NewConfigFromEnv with SESSION_KEY
func (suite *ConfigTestSuite) TestNewConfigFromEnv_WithSessionKey() {
	sessionKeyHex := "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210"
	os.Setenv("SESSION_KEY", sessionKeyHex)
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Len(suite.T(), config.SessionKey, 32)
}

// Test NewConfigFromEnv with invalid SESSION_KEY
func (suite *ConfigTestSuite) TestNewConfigFromEnv_InvalidSessionKey() {
	os.Setenv("SESSION_KEY", "invalid-hex-string")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), config)
}

// Test NewConfigFromEnv with TOKEN_EXPIRATION_HOURS
func (suite *ConfigTestSuite) TestNewConfigFromEnv_WithTokenExpiration() {
	os.Setenv("TOKEN_EXPIRATION_HOURS", "24")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Equal(suite.T(), 24*time.Hour, config.TokenExpiration)
}

// Test NewConfigFromEnv with invalid TOKEN_EXPIRATION_HOURS
func (suite *ConfigTestSuite) TestNewConfigFromEnv_InvalidTokenExpiration() {
	os.Setenv("TOKEN_EXPIRATION_HOURS", "not-a-number")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), config)
}

// Test NewConfigFromEnv with REFRESH_TOKEN_EXPIRATION_HOURS
func (suite *ConfigTestSuite) TestNewConfigFromEnv_WithRefreshTokenExpiration() {
	os.Setenv("REFRESH_TOKEN_EXPIRATION_HOURS", "168")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Equal(suite.T(), 168*time.Hour, config.RefreshTokenExp)
}

// Test NewConfigFromEnv with invalid REFRESH_TOKEN_EXPIRATION_HOURS
func (suite *ConfigTestSuite) TestNewConfigFromEnv_InvalidRefreshTokenExpiration() {
	os.Setenv("REFRESH_TOKEN_EXPIRATION_HOURS", "invalid")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), config)
}

// Test NewConfigFromEnv with BCRYPT_COST
func (suite *ConfigTestSuite) TestNewConfigFromEnv_WithBcryptCost() {
	os.Setenv("BCRYPT_COST", "12")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Equal(suite.T(), 12, config.BcryptCost)
}

// Test NewConfigFromEnv with invalid BCRYPT_COST
func (suite *ConfigTestSuite) TestNewConfigFromEnv_InvalidBcryptCost() {
	os.Setenv("BCRYPT_COST", "not-a-number")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), config)
}

// Test NewConfigFromEnv with ENABLE_RATE_LIMIT
func (suite *ConfigTestSuite) TestNewConfigFromEnv_WithRateLimit() {
	os.Setenv("ENABLE_RATE_LIMIT", "true")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.True(suite.T(), config.EnableRateLimit)
}

// Test NewConfigFromEnv with ENABLE_RATE_LIMIT false
func (suite *ConfigTestSuite) TestNewConfigFromEnv_WithRateLimitFalse() {
	os.Setenv("ENABLE_RATE_LIMIT", "false")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.False(suite.T(), config.EnableRateLimit)
}

// Test NewConfigFromEnv with MAX_LOGIN_ATTEMPTS
func (suite *ConfigTestSuite) TestNewConfigFromEnv_WithMaxLoginAttempts() {
	os.Setenv("MAX_LOGIN_ATTEMPTS", "5")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Equal(suite.T(), 5, config.MaxLoginAttempts)
}

// Test NewConfigFromEnv with invalid MAX_LOGIN_ATTEMPTS
func (suite *ConfigTestSuite) TestNewConfigFromEnv_InvalidMaxLoginAttempts() {
	os.Setenv("MAX_LOGIN_ATTEMPTS", "invalid")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), config)
}

// Test NewConfigFromEnv with RATE_LIMIT_WINDOW_MINUTES
func (suite *ConfigTestSuite) TestNewConfigFromEnv_WithRateLimitWindow() {
	os.Setenv("RATE_LIMIT_WINDOW_MINUTES", "15")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Equal(suite.T(), 15*time.Minute, config.RateLimitWindow)
}

// Test NewConfigFromEnv with invalid RATE_LIMIT_WINDOW_MINUTES
func (suite *ConfigTestSuite) TestNewConfigFromEnv_InvalidRateLimitWindow() {
	os.Setenv("RATE_LIMIT_WINDOW_MINUTES", "not-a-number")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), config)
}

// Test NewConfigFromEnv with RATE_LIMIT_BLOCK_MINUTES
func (suite *ConfigTestSuite) TestNewConfigFromEnv_WithRateLimitBlock() {
	os.Setenv("RATE_LIMIT_BLOCK_MINUTES", "30")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Equal(suite.T(), 30*time.Minute, config.RateLimitBlock)
}

// Test NewConfigFromEnv with invalid RATE_LIMIT_BLOCK_MINUTES
func (suite *ConfigTestSuite) TestNewConfigFromEnv_InvalidRateLimitBlock() {
	os.Setenv("RATE_LIMIT_BLOCK_MINUTES", "invalid")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), config)
}

// Test NewConfigFromEnv with all environment variables
func (suite *ConfigTestSuite) TestNewConfigFromEnv_AllEnvVars() {
	// Set all environment variables
	os.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	os.Setenv("SESSION_KEY", "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210")
	os.Setenv("TOKEN_EXPIRATION_HOURS", "24")
	os.Setenv("REFRESH_TOKEN_EXPIRATION_HOURS", "168")
	os.Setenv("BCRYPT_COST", "12")
	os.Setenv("ENABLE_RATE_LIMIT", "true")
	os.Setenv("MAX_LOGIN_ATTEMPTS", "5")
	os.Setenv("RATE_LIMIT_WINDOW_MINUTES", "15")
	os.Setenv("RATE_LIMIT_BLOCK_MINUTES", "30")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	
	// Verify all values are set correctly
	assert.Len(suite.T(), config.JWTSecret, 32)
	assert.Len(suite.T(), config.SessionKey, 32)
	assert.Equal(suite.T(), 24*time.Hour, config.TokenExpiration)
	assert.Equal(suite.T(), 168*time.Hour, config.RefreshTokenExp)
	assert.Equal(suite.T(), 12, config.BcryptCost)
	assert.True(suite.T(), config.EnableRateLimit)
	assert.Equal(suite.T(), 5, config.MaxLoginAttempts)
	assert.Equal(suite.T(), 15*time.Minute, config.RateLimitWindow)
	assert.Equal(suite.T(), 30*time.Minute, config.RateLimitBlock)
	assert.Equal(suite.T(), suite.logger, config.Logger)
}

// Test NewConfigFromEnv with zero values
func (suite *ConfigTestSuite) TestNewConfigFromEnv_ZeroValues() {
	os.Setenv("TOKEN_EXPIRATION_HOURS", "0")
	os.Setenv("REFRESH_TOKEN_EXPIRATION_HOURS", "0")
	os.Setenv("BCRYPT_COST", "0")
	os.Setenv("MAX_LOGIN_ATTEMPTS", "0")
	os.Setenv("RATE_LIMIT_WINDOW_MINUTES", "0")
	os.Setenv("RATE_LIMIT_BLOCK_MINUTES", "0")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Zero(suite.T(), config.TokenExpiration)
	assert.Zero(suite.T(), config.RefreshTokenExp)
	assert.Zero(suite.T(), config.BcryptCost)
	assert.Zero(suite.T(), config.MaxLoginAttempts)
	assert.Zero(suite.T(), config.RateLimitWindow)
	assert.Zero(suite.T(), config.RateLimitBlock)
}

// Test NewConfigFromEnv with negative values
func (suite *ConfigTestSuite) TestNewConfigFromEnv_NegativeValues() {
	os.Setenv("TOKEN_EXPIRATION_HOURS", "-1")
	os.Setenv("BCRYPT_COST", "-5")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Equal(suite.T(), -1*time.Hour, config.TokenExpiration)
	assert.Equal(suite.T(), -5, config.BcryptCost)
}

// Test NewConfigFromEnv with very large values
func (suite *ConfigTestSuite) TestNewConfigFromEnv_LargeValues() {
	os.Setenv("TOKEN_EXPIRATION_HOURS", "1000")
	os.Setenv("MAX_LOGIN_ATTEMPTS", "1000")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	// Values should be parsed correctly (even if impractical)
	assert.Equal(suite.T(), 1000*time.Hour, config.TokenExpiration)
	assert.Equal(suite.T(), 1000, config.MaxLoginAttempts)
}

// Test NewConfigFromEnv with empty string values
func (suite *ConfigTestSuite) TestNewConfigFromEnv_EmptyStringValues() {
	os.Setenv("JWT_SECRET", "")
	os.Setenv("SESSION_KEY", "")
	os.Setenv("TOKEN_EXPIRATION_HOURS", "")
	os.Setenv("ENABLE_RATE_LIMIT", "")
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Empty(suite.T(), config.JWTSecret)
	assert.Empty(suite.T(), config.SessionKey)
	assert.Zero(suite.T(), config.TokenExpiration)
	assert.False(suite.T(), config.EnableRateLimit)
}

// Test NewConfigFromEnv with nil logger
func (suite *ConfigTestSuite) TestNewConfigFromEnv_NilLogger() {
	config, err := NewConfigFromEnv(nil)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
	assert.Nil(suite.T(), config.Logger)
}

// Test edge case: odd-length hex string
func (suite *ConfigTestSuite) TestNewConfigFromEnv_OddLengthHex() {
	os.Setenv("JWT_SECRET", "123") // Odd length
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), config)
}

// Test edge case: hex string with invalid characters
func (suite *ConfigTestSuite) TestNewConfigFromEnv_InvalidHexChars() {
	os.Setenv("JWT_SECRET", "gghhiijj") // Invalid hex characters
	
	config, err := NewConfigFromEnv(suite.logger)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), config)
}

// Run the test suite
func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

// Additional unit tests for edge cases
func TestConfig_Struct(t *testing.T) {
	// Test that Config struct can be created manually
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	config := &Config{
		JWTSecret:        []byte("test-secret"),
		SessionKey:       []byte("test-session"),
		TokenExpiration:  time.Hour,
		RefreshTokenExp:  24 * time.Hour,
		BcryptCost:       10,
		Logger:           logger,
		EnableRateLimit:  true,
		MaxLoginAttempts: 3,
		RateLimitWindow:  time.Minute,
		RateLimitBlock:   time.Hour,
	}
	
	assert.NotNil(t, config)
	assert.Equal(t, []byte("test-secret"), config.JWTSecret)
	assert.Equal(t, []byte("test-session"), config.SessionKey)
	assert.Equal(t, time.Hour, config.TokenExpiration)
	assert.Equal(t, 24*time.Hour, config.RefreshTokenExp)
	assert.Equal(t, 10, config.BcryptCost)
	assert.Equal(t, logger, config.Logger)
	assert.True(t, config.EnableRateLimit)
	assert.Equal(t, 3, config.MaxLoginAttempts)
	assert.Equal(t, time.Minute, config.RateLimitWindow)
	assert.Equal(t, time.Hour, config.RateLimitBlock)
}

func TestConfig_Integration(t *testing.T) {
	// Test that config works with AuthService
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	
	config := &Config{
		JWTSecret:       []byte("test-secret-key-32-bytes-long!!"),
		SessionKey:      []byte("test-session-key-32-bytes-long!"),
		TokenExpiration: time.Hour,
		RefreshTokenExp: 24 * time.Hour,
		BcryptCost:      4, // Low cost for testing
		Logger:          logger,
		EnableRateLimit: false,
	}
	
	mockRepo := new(MockUserRepository)
	authService, err := NewAuthService(mockRepo, config)
	
	require.NoError(t, err)
	assert.NotNil(t, authService)
}