package auth

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	"github.com/mymindmap/api/models"
)

// MockUserRepository is a mock implementation of UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		user.ID = 1 // Mock setting ID
	}
	return args.Error(0)
}

// AuthServiceTestSuite defines the test suite for AuthService
type AuthServiceTestSuite struct {
	suite.Suite
	authService *AuthService
	mockRepo    *MockUserRepository
	config      *Config
	logger      *slog.Logger
}

// SetupTest runs before each test
func (suite *AuthServiceTestSuite) SetupTest() {
	suite.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Reduce noise in tests
	}))

	suite.config = &Config{
		JWTSecret:       []byte("test-secret-key-32-bytes-long!!"),
		SessionKey:      []byte("test-session-key-32-bytes-long!"),
		TokenExpiration: time.Hour,
		RefreshTokenExp: 24 * time.Hour,
		BcryptCost:      4, // Lower cost for faster tests
		Logger:          suite.logger,
		EnableRateLimit: false, // Disable for most tests
	}

	suite.mockRepo = new(MockUserRepository)
	
	var err error
	suite.authService, err = NewAuthService(suite.mockRepo, suite.config)
	require.NoError(suite.T(), err)
}

// TearDownTest runs after each test
func (suite *AuthServiceTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test NewAuthService
func (suite *AuthServiceTestSuite) TestNewAuthService_Success() {
	config := &Config{
		EnableRateLimit: true,
	}
	
	service, err := NewAuthService(suite.mockRepo, config)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), service)
	assert.NotNil(suite.T(), service.rateLimiter)
	assert.Equal(suite.T(), TokenExpirationTime, config.TokenExpiration)
	assert.Equal(suite.T(), BcryptCost, config.BcryptCost)
}

func (suite *AuthServiceTestSuite) TestNewAuthService_NilConfig() {
	_, err := NewAuthService(suite.mockRepo, nil)
	
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "config is required")
}

// Test RegisterUser
func (suite *AuthServiceTestSuite) TestRegisterUser_Success() {
	ctx := context.Background()
	req := &models.RegisterRequest{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "SecureP@ssw0rd123!",
	}

	// Mock that user doesn't exist
	suite.mockRepo.On("GetUserByEmail", ctx, "john.doe@example.com").Return(nil, nil)
	
	// Mock user creation
	suite.mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(nil)

	user, err := suite.authService.RegisterUser(ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "John Doe", user.Name)
	assert.Equal(suite.T(), "john.doe@example.com", user.Email)
	assert.Equal(suite.T(), RoleUser, user.Role)
	assert.Empty(suite.T(), user.Password) // Password should be cleared from response
}

func (suite *AuthServiceTestSuite) TestRegisterUser_UserExists() {
	ctx := context.Background()
	req := &models.RegisterRequest{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "SecureP@ssw0rd123!",
	}

	existingUser := &models.User{
		ID:    1,
		Email: "john.doe@example.com",
	}

	suite.mockRepo.On("GetUserByEmail", ctx, "john.doe@example.com").Return(existingUser, nil)

	user, err := suite.authService.RegisterUser(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), ErrUserExists, err)
}

func (suite *AuthServiceTestSuite) TestRegisterUser_InvalidEmail() {
	ctx := context.Background()
	req := &models.RegisterRequest{
		Name:     "John Doe",
		Email:    "invalid-email",
		Password: "SecureP@ssw0rd123!",
	}

	user, err := suite.authService.RegisterUser(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), ErrInvalidEmail, err)
}

func (suite *AuthServiceTestSuite) TestRegisterUser_WeakPassword() {
	ctx := context.Background()
	req := &models.RegisterRequest{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "weak",
	}

	user, err := suite.authService.RegisterUser(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "password must be at least")
}

// Test LoginUser
func (suite *AuthServiceTestSuite) TestLoginUser_Success() {
	ctx := context.Background()
	req := &models.LoginRequest{
		Email:    "john.doe@example.com",
		Password: "SecureP@ssw0rd123!",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("SecureP@ssw0rd123!"), suite.config.BcryptCost)
	user := &models.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: string(hashedPassword),
		Role:     RoleUser,
	}

	suite.mockRepo.On("GetUserByEmail", ctx, "john.doe@example.com").Return(user, nil)

	tokenPair, err := suite.authService.LoginUser(ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), tokenPair)
	assert.NotEmpty(suite.T(), tokenPair.AccessToken)
	assert.NotEmpty(suite.T(), tokenPair.RefreshToken)
	assert.Greater(suite.T(), tokenPair.ExpiresAt, time.Now().Unix())
}

func (suite *AuthServiceTestSuite) TestLoginUser_UserNotFound() {
	ctx := context.Background()
	req := &models.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "SecureP@ssw0rd123!",
	}

	suite.mockRepo.On("GetUserByEmail", ctx, "nonexistent@example.com").Return(nil, nil)

	tokenPair, err := suite.authService.LoginUser(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), tokenPair)
	assert.Equal(suite.T(), ErrInvalidCredentials, err)
}

func (suite *AuthServiceTestSuite) TestLoginUser_WrongPassword() {
	ctx := context.Background()
	req := &models.LoginRequest{
		Email:    "john.doe@example.com",
		Password: "WrongPassword123!",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("CorrectPassword123!"), suite.config.BcryptCost)
	user := &models.User{
		ID:       1,
		Email:    "john.doe@example.com",
		Password: string(hashedPassword),
		Role:     RoleUser,
	}

	suite.mockRepo.On("GetUserByEmail", ctx, "john.doe@example.com").Return(user, nil)

	tokenPair, err := suite.authService.LoginUser(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), tokenPair)
	assert.Equal(suite.T(), ErrInvalidCredentials, err)
}

// Test ValidateToken
func (suite *AuthServiceTestSuite) TestValidateToken_Success() {
	// Create a test user and token
	user := &models.User{
		ID:   1,
		Name: "John Doe",
		Email: "john.doe@example.com",
		Role: RoleUser,
	}

	token, err := suite.authService.createJWTToken(user, time.Hour)
	require.NoError(suite.T(), err)

	claims, err := suite.authService.ValidateToken(token)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), user.ID, claims.UserID)
	assert.Equal(suite.T(), user.Email, claims.Email)
	assert.Equal(suite.T(), user.Role, claims.Role)
}

func (suite *AuthServiceTestSuite) TestValidateToken_EmptyToken() {
	claims, err := suite.authService.ValidateToken("")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
	assert.Equal(suite.T(), ErrInvalidToken, err)
}

func (suite *AuthServiceTestSuite) TestValidateToken_InvalidToken() {
	claims, err := suite.authService.ValidateToken("invalid.jwt.token")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
}

// Test RefreshToken
func (suite *AuthServiceTestSuite) TestRefreshToken_Success() {
	user := &models.User{
		ID:   1,
		Name: "John Doe",
		Email: "john.doe@example.com",
		Role: RoleUser,
	}

	refreshToken, err := suite.authService.createJWTToken(user, 24*time.Hour)
	require.NoError(suite.T(), err)

	suite.mockRepo.On("GetUserByID", mock.Anything, 1).Return(user, nil)

	newTokenPair, err := suite.authService.RefreshToken(refreshToken)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), newTokenPair)
	assert.NotEmpty(suite.T(), newTokenPair.AccessToken)
	assert.NotEmpty(suite.T(), newTokenPair.RefreshToken)
}

func (suite *AuthServiceTestSuite) TestRefreshToken_InvalidToken() {
	newTokenPair, err := suite.authService.RefreshToken("invalid.token")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), newTokenPair)
}

// Test CheckPermission
func (suite *AuthServiceTestSuite) TestCheckPermission() {
	// Test user permissions
	assert.True(suite.T(), suite.authService.CheckPermission(RoleUser, ObjectPost, ActionRead))
	assert.True(suite.T(), suite.authService.CheckPermission(RoleUser, ObjectPost, ActionWrite))
	assert.False(suite.T(), suite.authService.CheckPermission(RoleUser, ObjectPost, ActionDelete))
	assert.False(suite.T(), suite.authService.CheckPermission(RoleUser, ObjectUser, ActionManage))

	// Test admin permissions
	assert.True(suite.T(), suite.authService.CheckPermission(RoleAdmin, ObjectPost, ActionRead))
	assert.True(suite.T(), suite.authService.CheckPermission(RoleAdmin, ObjectPost, ActionWrite))
	assert.True(suite.T(), suite.authService.CheckPermission(RoleAdmin, ObjectPost, ActionDelete))
	assert.True(suite.T(), suite.authService.CheckPermission(RoleAdmin, ObjectUser, ActionManage))
}

// Test CheckPermissionForUser
func (suite *AuthServiceTestSuite) TestCheckPermissionForUser() {
	email := "user@example.com"
	
	// Add user role
	err := suite.authService.AddRoleForUser(email, RoleUser)
	require.NoError(suite.T(), err)

	assert.True(suite.T(), suite.authService.CheckPermissionForUser(email, ObjectPost, ActionRead))
	assert.False(suite.T(), suite.authService.CheckPermissionForUser(email, ObjectPost, ActionDelete))
}

// Test role management
func (suite *AuthServiceTestSuite) TestAddRoleForUser_Success() {
	email := "user@example.com"
	
	err := suite.authService.AddRoleForUser(email, RoleAdmin)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), RoleAdmin, suite.authService.GetUserRole(email))
}

func (suite *AuthServiceTestSuite) TestAddRoleForUser_InvalidRole() {
	email := "user@example.com"
	
	err := suite.authService.AddRoleForUser(email, "invalid_role")
	
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid role")
}

func (suite *AuthServiceTestSuite) TestRemoveRoleForUser() {
	email := "user@example.com"
	
	// Add role first
	err := suite.authService.AddRoleForUser(email, RoleAdmin)
	require.NoError(suite.T(), err)
	
	// Remove role
	err = suite.authService.RemoveRoleForUser(email, RoleAdmin)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), RoleUser, suite.authService.GetUserRole(email)) // Should default to user
}

// Test validation methods
func (suite *AuthServiceTestSuite) TestValidateEmail() {
	testCases := []struct {
		email    string
		expected error
	}{
		{"valid@example.com", nil},
		{"user.name@domain.co.uk", nil},
		{"", errors.New("email is required")},
		{"invalid-email", ErrInvalidEmail},
		{"@domain.com", ErrInvalidEmail},
		{"user@", ErrInvalidEmail},
	}

	for _, tc := range testCases {
		err := suite.authService.validateEmail(tc.email)
		if tc.expected == nil {
			assert.NoError(suite.T(), err, "Email: %s", tc.email)
		} else {
			assert.Equal(suite.T(), tc.expected, err, "Email: %s", tc.email)
		}
	}
}

func (suite *AuthServiceTestSuite) TestValidatePassword() {
	testCases := []struct {
		password string
		valid    bool
	}{
		{"SecureP@ssw0rd123!", true},
		{"AnotherGood1!", true},
		{"weak", false},                    // Too short
		{"NoNumbers!", false},              // No digits
		{"nonumbers123", false},            // No uppercase
		{"NOLOWERCASE123!", false},         // No lowercase
		{"NoSpecialChars123", false},       // No special chars
		{"", false},                        // Empty
	}

	for _, tc := range testCases {
		err := suite.authService.validatePassword(tc.password)
		if tc.valid {
			assert.NoError(suite.T(), err, "Password: %s", tc.password)
		} else {
			assert.Error(suite.T(), err, "Password: %s", tc.password)
		}
	}
}

// Test with rate limiting enabled
func (suite *AuthServiceTestSuite) TestLoginUser_WithRateLimit() {
	// Create service with rate limiting
	config := &Config{
		JWTSecret:        []byte("test-secret-key-32-bytes-long!!"),
		SessionKey:       []byte("test-session-key-32-bytes-long!"),
		TokenExpiration:  time.Hour,
		RefreshTokenExp:  24 * time.Hour,
		BcryptCost:       4,
		Logger:           suite.logger,
		EnableRateLimit:  true,
		MaxLoginAttempts: 2,
		RateLimitWindow:  time.Minute,
		RateLimitBlock:   time.Minute,
	}

	authService, err := NewAuthService(suite.mockRepo, config)
	require.NoError(suite.T(), err)

	ctx := context.Background()
	req := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "WrongPassword123!",
	}

	// Mock failed login attempts
	suite.mockRepo.On("GetUserByEmail", ctx, "test@example.com").Return(nil, nil).Times(2)

	// First two attempts should return invalid credentials
	_, err1 := authService.LoginUser(ctx, req)
	assert.Equal(suite.T(), ErrInvalidCredentials, err1)

	_, err2 := authService.LoginUser(ctx, req)
	assert.Equal(suite.T(), ErrInvalidCredentials, err2)

	// Third attempt should be rate limited
	_, err3 := authService.LoginUser(ctx, req)
	assert.Equal(suite.T(), ErrTooManyAttempts, err3)
}

// Run the test suite
func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

// Additional unit tests not in the suite
func TestConstants(t *testing.T) {
	assert.Equal(t, "user", RoleUser)
	assert.Equal(t, "admin", RoleAdmin)
	assert.Equal(t, "author", RoleAuthor)
	assert.Equal(t, "post", ObjectPost)
	assert.Equal(t, "user", ObjectUser)
	assert.Equal(t, "read", ActionRead)
	assert.Equal(t, "write", ActionWrite)
	assert.Equal(t, "delete", ActionDelete)
	assert.Equal(t, "manage", ActionManage)
}

func TestErrors(t *testing.T) {
	assert.Equal(t, "user with this email already exists", ErrUserExists.Error())
	assert.Equal(t, "invalid email or password", ErrInvalidCredentials.Error())
	assert.Equal(t, "invalid token", ErrInvalidToken.Error())
	assert.Equal(t, "password does not meet security requirements", ErrWeakPassword.Error())
	assert.Equal(t, "invalid email format", ErrInvalidEmail.Error())
	assert.Equal(t, "token has expired", ErrTokenExpired.Error())
	assert.Equal(t, "permission denied", ErrPermissionDenied.Error())
	assert.Equal(t, "too many login attempts, please try again later", ErrTooManyAttempts.Error())
}