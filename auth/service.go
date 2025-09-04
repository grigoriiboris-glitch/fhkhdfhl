package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/mymindmap/api/models"
	"github.com/mymindmap/api/repository"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Constants for roles, objects, and actions
const (
	// Roles
	RoleUser   = "user"
	RoleAdmin  = "admin"
	RoleAuthor = "author"

	// Objects
	ObjectPost = "post"
	ObjectUser = "user"

	// Actions
	ActionRead   = "read"
	ActionWrite  = "write"
	ActionDelete = "delete"
	ActionManage = "manage"

	// Token settings
	TokenExpirationTime = 24 * time.Hour
	RefreshTokenExpTime = 7 * 24 * time.Hour

	// Password constraints
	MinPasswordLength = 8
	MaxPasswordLength = 128
	BcryptCost        = 12
)

var (
	ErrUserExists         = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrWeakPassword       = errors.New("password does not meet security requirements")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrTokenExpired       = errors.New("token has expired")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrTooManyAttempts    = errors.New("too many login attempts, please try again later")
)

// Config holds configuration for the auth service
type Config struct {
	JWTSecret        []byte
	SessionKey       []byte
	TokenExpiration  time.Duration
	RefreshTokenExp  time.Duration
	BcryptCost       int
	Logger           *slog.Logger
	EnableRateLimit  bool
	MaxLoginAttempts int
	RateLimitWindow  time.Duration
	RateLimitBlock   time.Duration
}

type AuthService struct {
	userRepo    *repository.UserRepository
	enforcer    *casbin.Enforcer
	config      *Config
	logger      *slog.Logger
	rateLimiter *RateLimiter
}

type Claims struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func NewAuthService(userRepo *repository.UserRepository, config *Config) (*AuthService, error) {
	if config == nil {
		return nil, errors.New("config is required")
	}

	// Set defaults if not provided
	if config.TokenExpiration == 0 {
		config.TokenExpiration = TokenExpirationTime
	}
	if config.RefreshTokenExp == 0 {
		config.RefreshTokenExp = RefreshTokenExpTime
	}
	if config.BcryptCost == 0 {
		config.BcryptCost = BcryptCost
	}
	if config.MaxLoginAttempts == 0 {
		config.MaxLoginAttempts = 5
	}
	if config.RateLimitWindow == 0 {
		config.RateLimitWindow = 15 * time.Minute
	}
	if config.RateLimitBlock == 0 {
		config.RateLimitBlock = 15 * time.Minute
	}

	// Generate secrets if not provided (for development only)
	if len(config.JWTSecret) == 0 {
		config.JWTSecret = make([]byte, 32)
		if _, err := rand.Read(config.JWTSecret); err != nil {
			return nil, fmt.Errorf("failed to generate JWT secret: %w", err)
		}
		if config.Logger != nil {
			config.Logger.Warn("JWT secret was auto-generated. Use environment variables in production")
		}
	}

	if len(config.SessionKey) == 0 {
		config.SessionKey = make([]byte, 32)
		if _, err := rand.Read(config.SessionKey); err != nil {
			return nil, fmt.Errorf("failed to generate session key: %w", err)
		}
		if config.Logger != nil {
			config.Logger.Warn("Session key was auto-generated. Use environment variables in production")
		}
	}

	// Create Casbin model
	m, err := createCasbinModel()
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin model: %w", err)
	}

	// Create enforcer
	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	service := &AuthService{
		userRepo: userRepo,
		enforcer: enforcer,
		config:   config,
		logger:   config.Logger,
	}

	// Initialize rate limiter if enabled
	if config.EnableRateLimit {
		service.rateLimiter = NewRateLimiter(
			config.MaxLoginAttempts,
			config.RateLimitWindow,
			config.RateLimitBlock,
		)
	}

	// Initialize policies
	if err := service.initializePolicies(); err != nil {
		return nil, fmt.Errorf("failed to initialize policies: %w", err)
	}

	return service, nil
}

func createCasbinModel() (model.Model, error) {
	return model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
}

func (s *AuthService) initializePolicies() error {
	policies := [][]string{
		// User permissions
		{RoleUser, ObjectPost, ActionRead},
		{RoleUser, ObjectPost, ActionWrite},
		
		// Author permissions
		{RoleAuthor, ObjectPost, ActionRead},
		{RoleAuthor, ObjectPost, ActionWrite},
		
		// Admin permissions
		{RoleAdmin, ObjectPost, ActionRead},
		{RoleAdmin, ObjectPost, ActionWrite},
		{RoleAdmin, ObjectPost, ActionDelete},
		{RoleAdmin, ObjectUser, ActionManage},
	}

	for _, policy := range policies {
		if _, err := s.enforcer.AddPolicy(policy[0], policy[1], policy[2]); err != nil {
			return fmt.Errorf("failed to add policy %v: %w", policy, err)
		}
	}

	return nil
}

func (s *AuthService) RegisterUser(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	if err := s.validateRegistrationRequest(req); err != nil {
		return nil, err
	}

	// Check if user exists
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logError("failed to check existing user", err, "email", req.Email)
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.config.BcryptCost)
	if err != nil {
		s.logError("failed to hash password", err, "email", req.Email)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Name:     strings.TrimSpace(req.Name),
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: string(hashedPassword),
		Role:     RoleUser,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		s.logError("failed to create user", err, "email", user.Email)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Add role in Casbin
	if _, err := s.enforcer.AddRoleForUser(user.Email, user.Role); err != nil {
		s.logError("failed to add role for user", err, "email", user.Email, "role", user.Role)
		return nil, fmt.Errorf("failed to add role for user: %w", err)
	}

	s.logInfo("user registered successfully", "email", user.Email, "role", user.Role)
	
	// Don't return password hash
	user.Password = ""
	return user, nil
}

func (s *AuthService) LoginUser(ctx context.Context, req *models.LoginRequest) (*TokenPair, error) {
	if err := s.validateLoginRequest(req); err != nil {
		return nil, err
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Check rate limiting
	if s.rateLimiter != nil && !s.rateLimiter.IsAllowed(email) {
		s.logInfo("login attempt blocked by rate limiter", "email", email)
		return nil, ErrTooManyAttempts
	}

	// Get user
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logError("failed to get user", err, "email", email)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		// Record failed attempt for rate limiting
		if s.rateLimiter != nil {
			s.rateLimiter.RecordAttempt(email)
		}
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logInfo("invalid password attempt", "email", email)
		// Record failed attempt for rate limiting
		if s.rateLimiter != nil {
			s.rateLimiter.RecordAttempt(email)
		}
		return nil, ErrInvalidCredentials
	}

	// Ensure role is in Casbin
	if _, err := s.enforcer.AddRoleForUser(user.Email, user.Role); err != nil {
		s.logError("failed to add role for user", err, "email", user.Email, "role", user.Role)
		return nil, fmt.Errorf("failed to add role for user: %w", err)
	}

	// Create token pair
	tokenPair, err := s.createTokenPair(user)
	if err != nil {
		s.logError("failed to create token pair", err, "email", user.Email)
		return nil, fmt.Errorf("failed to create token pair: %w", err)
	}

	// Reset rate limiting on successful login
	if s.rateLimiter != nil {
		s.rateLimiter.Reset(email)
	}

	s.logInfo("user logged in successfully", "email", user.Email)
	return tokenPair, nil
}

func (s *AuthService) createTokenPair(user *models.User) (*TokenPair, error) {
	// Create access token
	accessToken, err := s.createJWTToken(user, s.config.TokenExpiration)
	if err != nil {
		return nil, err
	}

	// Create refresh token
	refreshToken, err := s.createJWTToken(user, s.config.RefreshTokenExp)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.config.TokenExpiration).Unix(),
	}, nil
}

func (s *AuthService) createJWTToken(user *models.User, expiration time.Duration) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
			Issuer:    "mymindmap-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.config.JWTSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	if strings.TrimSpace(tokenString) == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.config.JWTSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *AuthService) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get fresh user data
	user, err := s.userRepo.GetUserByID(context.Background(), claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidToken
	}

	// Create new token pair
	return s.createTokenPair(user)
}

func (s *AuthService) CheckPermission(subject, object, action string) bool {
	allowed, err := s.enforcer.Enforce(subject, object, action)
	if err != nil {
		s.logError("permission check failed", err, "subject", subject, "object", object, "action", action)
		return false
	}
	return allowed
}

func (s *AuthService) CheckPermissionForUser(userEmail, object, action string) bool {
	roles, err := s.enforcer.GetRolesForUser(userEmail)
	if err != nil {
		s.logError("failed to get roles for user", err, "email", userEmail)
		return s.CheckPermission(RoleUser, object, action)
	}
	
	if len(roles) == 0 {
		return s.CheckPermission(RoleUser, object, action)
	}
	
	for _, role := range roles {
		if s.CheckPermission(role, object, action) {
			return true
		}
	}
	
	return false
}

func (s *AuthService) GetUserRole(email string) string {
	roles, err := s.enforcer.GetRolesForUser(email)
	if err != nil || len(roles) == 0 {
		return RoleUser
	}
	return roles[0]
}

func (s *AuthService) AddRoleForUser(email, role string) error {
	if !s.isValidRole(role) {
		return fmt.Errorf("invalid role: %s", role)
	}
	
	_, err := s.enforcer.AddRoleForUser(email, role)
	if err == nil {
		s.logInfo("role added for user", "email", email, "role", role)
	}
	return err
}

func (s *AuthService) RemoveRoleForUser(email, role string) error {
	_, err := s.enforcer.RemoveGroupingPolicy(email, role)
	if err == nil {
		s.logInfo("role removed for user", "email", email, "role", role)
	}
	return err
}

// Validation methods

func (s *AuthService) validateRegistrationRequest(req *models.RegisterRequest) error {
	if req == nil {
		return errors.New("registration request is required")
	}

	if err := s.validateEmail(req.Email); err != nil {
		return err
	}

	if err := s.validatePassword(req.Password); err != nil {
		return err
	}

	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}

	return nil
}

func (s *AuthService) validateLoginRequest(req *models.LoginRequest) error {
	if req == nil {
		return errors.New("login request is required")
	}

	if err := s.validateEmail(req.Email); err != nil {
		return err
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password is required")
	}

	return nil
}

func (s *AuthService) validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}

func (s *AuthService) validatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", MinPasswordLength)
	}

	if len(password) > MaxPasswordLength {
		return fmt.Errorf("password must not exceed %d characters", MaxPasswordLength)
	}

	// Check for at least one digit, one lowercase, one uppercase, and one special character
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`).MatchString(password)

	if !hasDigit || !hasLower || !hasUpper || !hasSpecial {
		return ErrWeakPassword
	}

	return nil
}

func (s *AuthService) isValidRole(role string) bool {
	validRoles := []string{RoleUser, RoleAdmin, RoleAuthor}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// Logging helpers

func (s *AuthService) logInfo(msg string, args ...any) {
	if s.logger != nil {
		s.logger.Info(msg, args...)
	}
}

func (s *AuthService) logError(msg string, err error, args ...any) {
	if s.logger != nil {
		allArgs := append(args, "error", err)
		s.logger.Error(msg, allArgs...)
	}
}