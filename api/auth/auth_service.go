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
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

<<<<<<< HEAD
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

=======
// UserRepositoryInterface определяет интерфейс для операций с пользователями в БД
// Это абстракция, позволяющая работать с разными базами данных
type UserRepositoryInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

// Константы для ролей, объектов и действий в системе прав доступа
const (
	// Роли пользователей
	RoleUser   = "user"    // Обычный пользователь
	RoleAdmin  = "admin"   // Администратор (полные права)
	RoleAuthor = "author"  // Автор контента

	// Объекты для контроля доступа
	ObjectPost = "post"    // Посты/записи
	ObjectUser = "user"    // Пользователи

	// Действия над объектами
	ActionRead   = "read"    // Просмотр
	ActionWrite  = "write"   // Создание/редактирование
	ActionDelete = "delete"  // Удаление
	ActionManage = "manage"  // Полное управление

	// Настройки токенов
	TokenExpirationTime = 24 * time.Hour      // Время жизни access токена
	RefreshTokenExpTime = 7 * 24 * time.Hour  // Время жизни refresh токена

	// Ограничения паролей
	MinPasswordLength = 8     // Минимальная длина пароля
	MaxPasswordLength = 128   // Максимальная длина пароля
	BcryptCost        = 12    // Сложность хеширования (12 - хороший баланс безопасности и производительности)
)

// Ошибки аутентификации и авторизации
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
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

<<<<<<< HEAD
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
=======
// Config содержит конфигурацию сервиса аутентификации
// Все настройки можно передать из внешнего config файла или переменных окружения
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

// AuthService - основной сервис аутентификации и авторизации
// Содержит бизнес-логику работы с пользователями, токенами и правами доступа
type AuthService struct {
	userRepo    UserRepositoryInterface // Репозиторий для работы с БД
	enforcer    *casbin.Enforcer        // Casbin enforcer для контроля доступа
	config      *Config                 // Конфигурация сервиса
	logger      *slog.Logger            // Логгер
	rateLimiter *RateLimiter            // Лимитер запросов (опционально)
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
}

// Claims - кастомные claims для JWT токена
// Содержат информацию о пользователе и стандартные JWT claims
type Claims struct {
	UserID int    `json:"user_id"`  // ID пользователя в БД
	Name   string `json:"name"`     // Имя пользователя
	Email  string `json:"email"`    // Email пользователя
	Role   string `json:"role"`     // Роль пользователя в системе
	jwt.RegisteredClaims            // Стандартные JWT claims (exp, iat, nbf, iss, etc.)
}

<<<<<<< HEAD
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
=======
// TokenPair - пара access и refresh токенов
// Используется для аутентификации и обновления сессии
type TokenPair struct {
	AccessToken  string `json:"access_token"`  // Короткоживущий токен для доступа к API
	RefreshToken string `json:"refresh_token"` // Долгоживущий токен для обновления access токена
	ExpiresAt    int64  `json:"expires_at"`    // Unix timestamp истечения access токена
}
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438

// NewAuthService создает новый экземпляр сервиса аутентификации
// Инициализирует все зависимости: Casbin, лимитер, настройки
func NewAuthService(userRepo UserRepositoryInterface, config *Config) (*AuthService, error) {
	if config == nil {
		return nil, errors.New("config is required")
	}

	// Установка значений по умолчанию, если не предоставлены
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

	// Генерация секретов, если не предоставлены (ТОЛЬКО для разработки!)
	if len(config.JWTSecret) == 0 {
		config.JWTSecret = make([]byte, 32)
		if _, err := rand.Read(config.JWTSecret); err != nil {
			return nil, fmt.Errorf("failed to generate JWT secret: %w", err)
		}
		if config.Logger != nil {
			config.Logger.Warn("JWT secret was auto-generated. Use environment variables in production")
		}
	}

<<<<<<< HEAD
[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
=======
	if len(config.SessionKey) == 0 {
		config.SessionKey = make([]byte, 32)
		if _, err := rand.Read(config.SessionKey); err != nil {
			return nil, fmt.Errorf("failed to generate session key: %w", err)
		}
		if config.Logger != nil {
			config.Logger.Warn("Session key was auto-generated. Use environment variables in production")
		}
	}

	// Создание Casbin модели для контроля доступа
	m, err := createCasbinModel()
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin model: %w", err)
	}

	// Создание Casbin enforcer - движка контроля доступа
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

	// Инициализация лимитера запросов, если включен
	if config.EnableRateLimit {
		service.rateLimiter = NewRateLimiter(
			config.MaxLoginAttempts,
			config.RateLimitWindow,
			config.RateLimitBlock,
		)
	}

	// Инициализация политик доступа
	if err := service.initializePolicies(); err != nil {
		return nil, fmt.Errorf("failed to initialize policies: %w", err)
	}

	return service, nil
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
}

// createCasbinModel создает модель Casbin для контроля доступа
// Используется RBAC (Role-Based Access Control) с субъект-объект-действие
func createCasbinModel() (model.Model, error) {
	return model.NewModelFromString(`
[request_definition]
r = sub, obj, act  # Запрос: кто (роль), что (объект), какое действие

[policy_definition]
p = sub, obj, act  # Политика: для какой роли, на какой объект, какое действие разрешено

[role_definition]
g = _, _           # Назначение ролей пользователям (user -> role)

[policy_effect]
e = some(where (p.eft == allow))  # Эффект: разрешить если хотя бы одна политика позволяет

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act  # Совпадение: пользователь имеет роль И объект совпадает И действие совпадает
`)
}

// initializePolicies инициализирует базовые политики доступа
// Определяет какие роли могут выполнять какие действия над какими объектами
func (s *AuthService) initializePolicies() error {
	policies := [][]string{
<<<<<<< HEAD
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

=======
		// Права обычного пользователя
		{RoleUser, ObjectPost, ActionRead},   // user может читать посты
		{RoleUser, ObjectPost, ActionWrite},  // user может писать посты
		
		// Права автора (наследует права user + дополнительные)
		{RoleAuthor, ObjectPost, ActionRead},  // author может читать посты
		{RoleAuthor, ObjectPost, ActionWrite}, // author может писать посты
		
		// Права администратора (полные права)
		{RoleAdmin, ObjectPost, ActionRead},    // admin может читать посты
		{RoleAdmin, ObjectPost, ActionWrite},   // admin может писать посты
		{RoleAdmin, ObjectPost, ActionDelete},  // admin может удалять посты
		{RoleAdmin, ObjectUser, ActionManage},  // admin может управлять пользователями
	}

	// Добавление всех политик в enforcer
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	for _, policy := range policies {
		if _, err := s.enforcer.AddPolicy(policy[0], policy[1], policy[2]); err != nil {
			return fmt.Errorf("failed to add policy %v: %w", policy, err)
		}
	}

	return nil
}

// RegisterUser регистрирует нового пользователя в системе
// Валидирует данные, хеширует пароль, создает запись в БД и назначает роль
func (s *AuthService) RegisterUser(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	if err := s.validateRegistrationRequest(req); err != nil {
		return nil, err
	}

<<<<<<< HEAD
	// Check if user exists
=======
	// Проверка существования пользователя с таким email
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logError("failed to check existing user", err, "email", req.Email)
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

<<<<<<< HEAD
	// Hash password
=======
	// Хеширование пароля с bcrypt
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.config.BcryptCost)
	if err != nil {
		s.logError("failed to hash password", err, "email", req.Email)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

<<<<<<< HEAD
	// Create user
=======
	// Создание объекта пользователя
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	user := &models.User{
		Name:     strings.TrimSpace(req.Name),
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: string(hashedPassword),
<<<<<<< HEAD
		Role:     RoleUser,
=======
		Role:     RoleUser, // По умолчанию обычный пользователь
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	}

	// Сохранение пользователя в БД
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		s.logError("failed to create user", err, "email", user.Email)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

<<<<<<< HEAD
	// Add role in Casbin
=======
	// Назначение роли в Casbin
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	if _, err := s.enforcer.AddRoleForUser(user.Email, user.Role); err != nil {
		s.logError("failed to add role for user", err, "email", user.Email, "role", user.Role)
		return nil, fmt.Errorf("failed to add role for user: %w", err)
	}

	s.logInfo("user registered successfully", "email", user.Email, "role", user.Role)
	
<<<<<<< HEAD
	// Don't return password hash
=======
	// Не возвращаем хеш пароля в ответе
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	user.Password = ""
	return user, nil
}

<<<<<<< HEAD
=======
// LoginUser аутентифицирует пользователя и выдает токены
// Проверяет учетные данные, лимиты запросов и создает JWT токены
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
func (s *AuthService) LoginUser(ctx context.Context, req *models.LoginRequest) (*TokenPair, error) {
	if err := s.validateLoginRequest(req); err != nil {
		return nil, err
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

<<<<<<< HEAD
	// Check rate limiting
=======
	// Проверка лимита запросов (защита от брутфорса)
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	if s.rateLimiter != nil && !s.rateLimiter.IsAllowed(email) {
		s.logInfo("login attempt blocked by rate limiter", "email", email)
		return nil, ErrTooManyAttempts
	}

<<<<<<< HEAD
	// Get user
=======
	// Получение пользователя из БД
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logError("failed to get user", err, "email", email)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
<<<<<<< HEAD
		// Record failed attempt for rate limiting
=======
		// Запись неудачной попытки для лимитера
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
		if s.rateLimiter != nil {
			s.rateLimiter.RecordAttempt(email)
		}
		return nil, ErrInvalidCredentials
	}

<<<<<<< HEAD
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logInfo("invalid password attempt", "email", email)
		// Record failed attempt for rate limiting
=======
	// Проверка пароля с bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logInfo("invalid password attempt", "email", email)
		// Запись неудачной попытки
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
		if s.rateLimiter != nil {
			s.rateLimiter.RecordAttempt(email)
		}
		return nil, ErrInvalidCredentials
	}

<<<<<<< HEAD
	// Ensure role is in Casbin
=======
	// Убеждаемся что роль назначена в Casbin
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	if _, err := s.enforcer.AddRoleForUser(user.Email, user.Role); err != nil {
		s.logError("failed to add role for user", err, "email", user.Email, "role", user.Role)
		return nil, fmt.Errorf("failed to add role for user: %w", err)
	}

<<<<<<< HEAD
	// Create token pair
=======
	// Создание пары токенов (access + refresh)
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	tokenPair, err := s.createTokenPair(user)
	if err != nil {
		s.logError("failed to create token pair", err, "email", user.Email)
		return nil, fmt.Errorf("failed to create token pair: %w", err)
	}

<<<<<<< HEAD
	// Reset rate limiting on successful login
=======
	// Сброс лимитера при успешном входе
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	if s.rateLimiter != nil {
		s.rateLimiter.Reset(email)
	}

	s.logInfo("user logged in successfully", "email", user.Email)
	return tokenPair, nil
}

<<<<<<< HEAD
func (s *AuthService) createTokenPair(user *models.User) (*TokenPair, error) {
	// Create access token
=======
// createTokenPair создает пару access и refresh токенов
// Access токен - короткоживущий, refresh - долгоживущий для обновления
func (s *AuthService) createTokenPair(user *models.User) (*TokenPair, error) {
	// Создание access токена
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	accessToken, err := s.createJWTToken(user, s.config.TokenExpiration)
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	// Create refresh token
=======
	// Создание refresh токена
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
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

<<<<<<< HEAD
=======
// createJWTToken создает JWT токен с указанным временем жизни
// Содержит claims с информацией о пользователе
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
func (s *AuthService) createJWTToken(user *models.User, expiration time.Duration) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
<<<<<<< HEAD
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
			Issuer:    "mymindmap-api",
=======
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)), // Время истечения
			IssuedAt:  jwt.NewNumericDate(time.Now()),                 // Время создания
			NotBefore: jwt.NewNumericDate(time.Now()),                 // Не действует до
			Subject:   fmt.Sprintf("%d", user.ID),                     // ID пользователя как subject
			Issuer:    "mymindmap-api",                                // Идентификатор издателя
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
<<<<<<< HEAD
	return token.SignedString(s.config.JWTSecret)
=======
	return token.SignedString(s.config.JWTSecret) // Подпись токена секретным ключом
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
}

// ValidateToken проверяет валидность JWT токена и возвращает claims
// Используется в middleware для аутентификации запросов
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	if strings.TrimSpace(tokenString) == "" {
		return nil, ErrInvalidToken
	}

<<<<<<< HEAD
=======
	// Парсинг токена с проверкой подписи
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
<<<<<<< HEAD
		return s.config.JWTSecret, nil
=======
		return s.config.JWTSecret, nil // Возвращаем секрет для проверки подписи
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

<<<<<<< HEAD
=======
	// Извлечение claims и проверка валидности токена
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
<<<<<<< HEAD
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
=======
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
}

// RefreshToken обновляет access токен с помощью refresh токена
// Проверяет refresh токен и выдает новую пару токенов
func (s *AuthService) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	// Валидация refresh токена
	claims, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Получение актуальных данных пользователя из БД
	user, err := s.userRepo.GetUserByID(context.Background(), claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidToken
	}

	// Создание новой пары токенов
	return s.createTokenPair(user)
}

// CheckPermission проверяет разрешение для конкретной роли
// sub - роль, obj - объект, act - действие
func (s *AuthService) CheckPermission(subject, object, action string) bool {
	allowed, err := s.enforcer.Enforce(subject, object, action)
	if err != nil {
		s.logError("permission check failed", err, "subject", subject, "object", object, "action", action)
		return false
	}
	return allowed
}

<<<<<<< HEAD
func (s *AuthService) CheckPermissionForUser(userEmail, object, action string) bool {
	roles, err := s.enforcer.GetRolesForUser(userEmail)
	if err != nil {
		s.logError("failed to get roles for user", err, "email", userEmail)
		return s.CheckPermission(RoleUser, object, action)
	}
	
=======
// CheckPermissionForUser проверяет разрешение для конкретного пользователя
// Определяет роли пользователя и проверяет права для каждой роли
func (s *AuthService) CheckPermissionForUser(userEmail, object, action string) bool {
	// Получение всех ролей пользователя
	roles, err := s.enforcer.GetRolesForUser(userEmail)
	if err != nil {
		s.logError("failed to get roles for user", err, "email", userEmail)
		return s.CheckPermission(RoleUser, object, action) // Fallback к user роли
	}
	
	// Если ролей нет - используем user роль
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
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

// GetUserRole возвращает основную роль пользователя
// Если ролей несколько - возвращает первую
func (s *AuthService) GetUserRole(email string) string {
	roles, err := s.enforcer.GetRolesForUser(email)
	if err != nil || len(roles) == 0 {
<<<<<<< HEAD
		return RoleUser
=======
		return RoleUser // Роль по умолчанию
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	}
	return roles[0]
}

// AddRoleForUser добавляет роль пользователю в системе прав доступа
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

// RemoveRoleForUser удаляет роль у пользователя
func (s *AuthService) RemoveRoleForUser(email, role string) error {
	_, err := s.enforcer.RemoveGroupingPolicy(email, role)
	if err == nil {
		s.logInfo("role removed for user", "email", email, "role", role)
	}
	return err
}

<<<<<<< HEAD
// Validation methods

=======
// Методы валидации

// validateRegistrationRequest валидирует данные для регистрации
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
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

<<<<<<< HEAD
=======
// validateLoginRequest валидирует данные для входа
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
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

<<<<<<< HEAD
=======
// validateEmail проверяет валидность email адреса
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
func (s *AuthService) validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}

<<<<<<< HEAD
=======
	// Регулярное выражение для проверки email
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}

<<<<<<< HEAD
=======
// validatePassword проверяет сложность пароля
// Требования: минимум 8 символов, цифры, буквы в разных регистрах, спецсимволы
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
func (s *AuthService) validatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", MinPasswordLength)
	}

	if len(password) > MaxPasswordLength {
		return fmt.Errorf("password must not exceed %d characters", MaxPasswordLength)
	}

<<<<<<< HEAD
	// Check for at least one digit, one lowercase, one uppercase, and one special character
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`).MatchString(password)
=======
	// Проверка сложности пароля
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)           // Есть цифры
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)           // Есть строчные буквы
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)           // Есть заглавные буквы
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`).MatchString(password) // Есть спецсимволы
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438

	if !hasDigit || !hasLower || !hasUpper || !hasSpecial {
		return ErrWeakPassword
	}

	return nil
}

<<<<<<< HEAD
=======
// isValidRole проверяет что роль является допустимой
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
func (s *AuthService) isValidRole(role string) bool {
	validRoles := []string{RoleUser, RoleAdmin, RoleAuthor}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

<<<<<<< HEAD
// Logging helpers

=======
// Вспомогательные методы логирования

// logInfo логирует информационное сообщение
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
func (s *AuthService) logInfo(msg string, args ...any) {
	if s.logger != nil {
		s.logger.Info(msg, args...)
	}
}

<<<<<<< HEAD
=======
// logError логирует сообщение об ошибке
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
func (s *AuthService) logError(msg string, err error, args ...any) {
	if s.logger != nil {
		allArgs := append(args, "error", err)
		s.logger.Error(msg, allArgs...)
	}
<<<<<<< HEAD
} 
=======
}
>>>>>>> cursor/organize-and-refactor-auth-code-into-api-auth-4438
