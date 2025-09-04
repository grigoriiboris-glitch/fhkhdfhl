package auth

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/mymindmap/api/models"
)

// MiddlewareTestSuite defines the test suite for middleware
type MiddlewareTestSuite struct {
	suite.Suite
	authService *AuthService
	mockRepo    *MockUserRepository
	config      *Config
}

// SetupTest runs before each test
func (suite *MiddlewareTestSuite) SetupTest() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	suite.config = &Config{
		JWTSecret:       []byte("test-secret-key-32-bytes-long!!"),
		SessionKey:      []byte("test-session-key-32-bytes-long!"),
		TokenExpiration: time.Hour,
		RefreshTokenExp: 24 * time.Hour,
		BcryptCost:      4,
		Logger:          logger,
		EnableRateLimit: false,
	}

	suite.mockRepo = new(MockUserRepository)
	
	var err error
	suite.authService, err = NewAuthService(suite.mockRepo, suite.config)
	require.NoError(suite.T(), err)
}

// TearDownTest runs after each test
func (suite *MiddlewareTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test AuthMiddleware with valid token in Authorization header
func (suite *MiddlewareTestSuite) TestAuthMiddleware_ValidTokenInHeader() {
	// Create a test user and token
	user := &models.User{
		ID:   1,
		Name: "John Doe",
		Email: "john.doe@example.com",
		Role: RoleUser,
	}

	token, err := suite.authService.createJWTToken(user, time.Hour)
	require.NoError(suite.T(), err)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserFromContext(r.Context())
		assert.NotNil(suite.T(), claims)
		assert.Equal(suite.T(), user.ID, claims.UserID)
		assert.Equal(suite.T(), user.Email, claims.Email)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with middleware
	handler := suite.authService.AuthMiddleware(testHandler)

	// Create request with Authorization header
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "success", w.Body.String())
}

// Test AuthMiddleware with valid token in cookie
func (suite *MiddlewareTestSuite) TestAuthMiddleware_ValidTokenInCookie() {
	// Create a test user and token
	user := &models.User{
		ID:   1,
		Name: "John Doe",
		Email: "john.doe@example.com",
		Role: RoleUser,
	}

	token, err := suite.authService.createJWTToken(user, time.Hour)
	require.NoError(suite.T(), err)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserFromContext(r.Context())
		assert.NotNil(suite.T(), claims)
		assert.Equal(suite.T(), user.ID, claims.UserID)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with middleware
	handler := suite.authService.AuthMiddleware(testHandler)

	// Create request with cookie
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth_token",
		Value: token,
	})
	
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "success", w.Body.String())
}

// Test AuthMiddleware with no token
func (suite *MiddlewareTestSuite) TestAuthMiddleware_NoToken() {
	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserFromContext(r.Context())
		assert.Nil(suite.T(), claims)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("no auth"))
	})

	// Wrap with middleware
	handler := suite.authService.AuthMiddleware(testHandler)

	// Create request without token
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "no auth", w.Body.String())
}

// Test AuthMiddleware with invalid token
func (suite *MiddlewareTestSuite) TestAuthMiddleware_InvalidToken() {
	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserFromContext(r.Context())
		assert.Nil(suite.T(), claims)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid token"))
	})

	// Wrap with middleware
	handler := suite.authService.AuthMiddleware(testHandler)

	// Create request with invalid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")
	
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "invalid token", w.Body.String())
}

// Test AuthMiddleware with expired token
func (suite *MiddlewareTestSuite) TestAuthMiddleware_ExpiredToken() {
	// Create a test user and expired token
	user := &models.User{
		ID:   1,
		Name: "John Doe",
		Email: "john.doe@example.com",
		Role: RoleUser,
	}

	// Create token with negative expiration (expired)
	token, err := suite.authService.createJWTToken(user, -time.Hour)
	require.NoError(suite.T(), err)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserFromContext(r.Context())
		assert.Nil(suite.T(), claims)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("expired token"))
	})

	// Wrap with middleware
	handler := suite.authService.AuthMiddleware(testHandler)

	// Create request with expired token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "expired token", w.Body.String())
}

// Test RequireAuth middleware with valid token
func (suite *MiddlewareTestSuite) TestRequireAuth_WithValidToken() {
	// Create a test user and token
	user := &models.User{
		ID:   1,
		Name: "John Doe",
		Email: "john.doe@example.com",
		Role: RoleUser,
	}

	token, err := suite.authService.createJWTToken(user, time.Hour)
	require.NoError(suite.T(), err)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authorized"))
	})

	// Wrap with both middlewares
	handler := suite.authService.AuthMiddleware(suite.authService.RequireAuth(testHandler))

	// Create request with token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "authorized", w.Body.String())
}

// Test RequireAuth middleware without token
func (suite *MiddlewareTestSuite) TestRequireAuth_WithoutToken() {
	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("should not reach here"))
	})

	// Wrap with RequireAuth middleware
	handler := suite.authService.RequireAuth(testHandler)

	// Create request without token
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Unauthorized")
}

// Test RequirePermission middleware
func (suite *MiddlewareTestSuite) TestRequirePermission_WithValidPermission() {
	// Create a test user and token
	user := &models.User{
		ID:   1,
		Name: "John Doe",
		Email: "john.doe@example.com",
		Role: RoleUser,
	}

	token, err := suite.authService.createJWTToken(user, time.Hour)
	require.NoError(suite.T(), err)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("permission granted"))
	})

	// Wrap with middlewares (user can read posts)
	handler := suite.authService.AuthMiddleware(
		suite.authService.RequirePermission(ObjectPost, ActionRead)(testHandler),
	)

	// Create request with token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "permission granted", w.Body.String())
}

// Test RequirePermission middleware without permission
func (suite *MiddlewareTestSuite) TestRequirePermission_WithoutPermission() {
	// Create a test user and token
	user := &models.User{
		ID:   1,
		Name: "John Doe",
		Email: "john.doe@example.com",
		Role: RoleUser,
	}

	token, err := suite.authService.createJWTToken(user, time.Hour)
	require.NoError(suite.T(), err)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("should not reach here"))
	})

	// Wrap with middlewares (user cannot delete posts)
	handler := suite.authService.AuthMiddleware(
		suite.authService.RequirePermission(ObjectPost, ActionDelete)(testHandler),
	)

	// Create request with token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Forbidden")
}

// Test RequirePermission middleware without authentication
func (suite *MiddlewareTestSuite) TestRequirePermission_WithoutAuth() {
	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("should not reach here"))
	})

	// Wrap with permission middleware only (no auth middleware)
	handler := suite.authService.RequirePermission(ObjectPost, ActionRead)(testHandler)

	// Create request without token
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Unauthorized")
}

// Test GetUserFromContext
func (suite *MiddlewareTestSuite) TestGetUserFromContext() {
	// Test with no user in context
	ctx := context.Background()
	user := GetUserFromContext(ctx)
	assert.Nil(suite.T(), user)

	// Test with user in context
	claims := &Claims{
		UserID: 1,
		Name:   "John Doe",
		Email:  "john.doe@example.com",
		Role:   RoleUser,
	}
	
	ctx = context.WithValue(ctx, UserContextKey, claims)
	user = GetUserFromContext(ctx)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), claims.UserID, user.UserID)
	assert.Equal(suite.T(), claims.Email, user.Email)

	// Test with wrong type in context
	ctx = context.WithValue(ctx, UserContextKey, "wrong type")
	user = GetUserFromContext(ctx)
	assert.Nil(suite.T(), user)
}

// Test SetAuthCookie
func (suite *MiddlewareTestSuite) TestSetAuthCookie() {
	tokenPair := &TokenPair{
		AccessToken:  "access.token.here",
		RefreshToken: "refresh.token.here",
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
	}

	w := httptest.NewRecorder()
	suite.authService.SetAuthCookie(w, tokenPair)

	cookies := w.Result().Cookies()
	assert.Len(suite.T(), cookies, 2)

	// Check access token cookie
	var accessCookie, refreshCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			accessCookie = cookie
		} else if cookie.Name == "refresh_token" {
			refreshCookie = cookie
		}
	}

	assert.NotNil(suite.T(), accessCookie)
	assert.Equal(suite.T(), "access.token.here", accessCookie.Value)
	assert.Equal(suite.T(), "/", accessCookie.Path)
	assert.True(suite.T(), accessCookie.HttpOnly)
	assert.Equal(suite.T(), int(suite.config.TokenExpiration.Seconds()), accessCookie.MaxAge)

	assert.NotNil(suite.T(), refreshCookie)
	assert.Equal(suite.T(), "refresh.token.here", refreshCookie.Value)
	assert.Equal(suite.T(), "/", refreshCookie.Path)
	assert.True(suite.T(), refreshCookie.HttpOnly)
	assert.Equal(suite.T(), int(suite.config.RefreshTokenExp.Seconds()), refreshCookie.MaxAge)
}

// Test ClearAuthCookie
func (suite *MiddlewareTestSuite) TestClearAuthCookie() {
	w := httptest.NewRecorder()
	suite.authService.ClearAuthCookie(w)

	cookies := w.Result().Cookies()
	assert.Len(suite.T(), cookies, 2)

	// Check that both cookies are cleared
	var accessCookie, refreshCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			accessCookie = cookie
		} else if cookie.Name == "refresh_token" {
			refreshCookie = cookie
		}
	}

	assert.NotNil(suite.T(), accessCookie)
	assert.Equal(suite.T(), "", accessCookie.Value)
	assert.Equal(suite.T(), -1, accessCookie.MaxAge)

	assert.NotNil(suite.T(), refreshCookie)
	assert.Equal(suite.T(), "", refreshCookie.Value)
	assert.Equal(suite.T(), -1, refreshCookie.MaxAge)
}

// Test middleware chaining
func (suite *MiddlewareTestSuite) TestMiddlewareChaining() {
	// Create admin user and token
	user := &models.User{
		ID:   1,
		Name: "Admin User",
		Email: "admin@example.com",
		Role: RoleAdmin,
	}

	token, err := suite.authService.createJWTToken(user, time.Hour)
	require.NoError(suite.T(), err)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserFromContext(r.Context())
		assert.NotNil(suite.T(), claims)
		assert.Equal(suite.T(), RoleAdmin, claims.Role)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("admin access granted"))
	})

	// Chain multiple middlewares: Auth -> RequireAuth -> RequirePermission
	handler := suite.authService.AuthMiddleware(
		suite.authService.RequireAuth(
			suite.authService.RequirePermission(ObjectUser, ActionManage)(testHandler),
		),
	)

	// Create request with admin token
	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "admin access granted", w.Body.String())
}

// Run the test suite
func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

// Additional unit tests
func TestContextKey(t *testing.T) {
	assert.Equal(t, contextKey("user"), UserContextKey)
}