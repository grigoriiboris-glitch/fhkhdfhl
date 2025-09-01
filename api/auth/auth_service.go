package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/mymindmap/api/models"
	"github.com/mymindmap/api/repository"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	enforcer    *casbin.Enforcer
	jwtSecret   []byte
	sessionKey  []byte
}

type Claims struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo *repository.UserRepository) (*AuthService, error) {
	// Создаем простую модель RBAC для Casbin
	m, err := model.NewModelFromString(`
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
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin model: %w", err)
	}

	// Создаем enforcer с пустой политикой (будет заполнена позже)
	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// Генерируем секретные ключи
	jwtSecret := make([]byte, 32)
	if _, err := rand.Read(jwtSecret); err != nil {
		return nil, fmt.Errorf("failed to generate JWT secret: %w", err)
	}

	sessionKey := make([]byte, 32)
	if _, err := rand.Read(sessionKey); err != nil {
		return nil, fmt.Errorf("failed to generate session key: %w", err)
	}

	service := &AuthService{
		userRepo:   userRepo,
		enforcer:   enforcer,
		jwtSecret:  jwtSecret,
		sessionKey: sessionKey,
	}

	// Инициализируем базовые политики
	if err := service.initializePolicies(); err != nil {
		return nil, fmt.Errorf("failed to initialize policies: %w", err)
	}

	return service, nil
}

func (s *AuthService) initializePolicies() error {
	// Добавляем базовые политики
	// Пользователи могут читать посты
	_, err := s.enforcer.AddPolicy("user", "post", "read")
	if err != nil {
		return err
	}

	// Пользователи могут создавать посты
	_, err = s.enforcer.AddPolicy("user", "post", "write")
	if err != nil {
		return err
	}

	// Администраторы могут делать все
	_, err = s.enforcer.AddPolicy("admin", "post", "read")
	if err != nil {
		return err
	}
	_, err = s.enforcer.AddPolicy("admin", "post", "write")
	if err != nil {
		return err
	}
	_, err = s.enforcer.AddPolicy("admin", "post", "delete")
	if err != nil {
		return err
	}
	_, err = s.enforcer.AddPolicy("admin", "user", "manage")
	if err != nil {
		return err
	}

	// Авторы могут писать и редактировать свои посты
	_, err = s.enforcer.AddPolicy("author", "post", "read")
	if err != nil {
		return err
	}
	_, err = s.enforcer.AddPolicy("author", "post", "write")
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RegisterUser(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаем пользователя
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user", // По умолчанию обычный пользователь
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Добавляем роль пользователя в Casbin
	_, err = s.enforcer.AddRoleForUser(user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to add role for user: %w", err)
	}

	return user, nil
}

func (s *AuthService) LoginUser(ctx context.Context, req *models.LoginRequest) (string, error) {
	// Получаем пользователя по email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Убеждаемся, что роль пользователя добавлена в Casbin
	_, err = s.enforcer.AddRoleForUser(user.Email, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to add role for user: %w", err)
	}

	// Создаем JWT токен
	token, err := s.createJWTToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
}

func (s *AuthService) createJWTToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) CheckPermission(subject, object, action string) bool {
	// Проверяем права напрямую по роли
	allowed, err := s.enforcer.Enforce(subject, object, action)
	if err != nil {
		return false
	}
	return allowed
}

// CheckPermissionForUser проверяет права для конкретного пользователя
func (s *AuthService) CheckPermissionForUser(userEmail, object, action string) bool {
	// Получаем роль пользователя
	roles, err := s.enforcer.GetRolesForUser(userEmail)
	if err != nil || len(roles) == 0 {
		// Если роль не найдена, используем "user" по умолчанию
		return s.CheckPermission("user", object, action)
	}
	
	// Проверяем права для каждой роли пользователя
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
		return "user"
	}
	return roles[0]
}

func (s *AuthService) AddRoleForUser(email, role string) error {
	_, err := s.enforcer.AddRoleForUser(email, role)
	return err
}

func (s *AuthService) RemoveRoleForUser(email, role string) error {
	_, err := s.enforcer.RemoveGroupingPolicy(email, role)
	return err
} 