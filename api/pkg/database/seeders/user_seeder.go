// pkg/database/seeders/user_seeder.go
package seeders

import (
	"fmt"
	"log"

	"github.com/mymindmap/api/models"
	"github.com/mymindmap/api/pkg/repositories"
)

// UserSeeder сидер для пользователей
type UserSeeder struct {
	BaseSeeder
	userRepo *repositories.UserRepository
}

// NewUserSeeder создает новый сидер для пользователей
func NewUserSeeder() *UserSeeder {
	baseSeeder := NewBaseSeeder()
	return &UserSeeder{
		BaseSeeder: *baseSeeder,
		userRepo:   repositories.NewUserRepository(baseSeeder.DB),
	}
}

// Run запускает сидер
func (s *UserSeeder) Run() error {
	log.Println("Запуск сидера пользователей...")

	// Очищаем таблицу перед заполнением
	if err := s.userRepo.Truncate(); err != nil {
		return fmt.Errorf("ошибка очистки таблицы: %v", err)
	}

	// Создаем тестовых пользователей
	users := []models.User{
		{
			Name:     "Администратор",
			Email:    "admin@example.com",
			Password: "password123",
			Role:     "admin",
		},
		{
			Name:     "Обычный пользователь",
			Email:    "user@example.com",
			Password: "password123",
			Role:     "user",
		},
		{
			Name:     "Тестовый пользователь",
			Email:    "test@example.com",
			Password: "password123",
			Role:     "user",
		},
	}

	for _, user := range users {
		if err := s.userRepo.Create(&user); err != nil {
			return fmt.Errorf("ошибка создания пользователя %s: %v", user.Email, err)
		}
		log.Printf("Создан пользователь: %s (%s)", user.Name, user.Email)
	}

	log.Println("Сидер пользователей завершен успешно!")
	return nil
}