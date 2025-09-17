package models

import (
	"time"
	"github.com/mymindmap/api/pkg/relations"
)

// теги для парсинга полей
type User struct {
	relations.BaseModel
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"` // Не отправляем пароль в JSON
	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	MindMaps  []*MindMap
	Posts  []*Post
}

func (u *User) GetID() interface{} {
    return u.ID
}

func (u *User) GetTableName() string { return "users" }

func (u *User) GetRelations() map[string]relations.Relation {
	return map[string]relations.Relation{
		"mindMaps": relations.NewHasMany(
			func() relations.Model { return &MindMap{} },
			relations.RelationConfig{
				ForeignKey: "user_id",
				LocalKey:   "id",
			},
		),
		"posts": relations.NewHasMany(
			func() relations.Model { return &Post{} },
			relations.RelationConfig{
				ForeignKey: "user_id",
				LocalKey:   "id",
			},
		),
	}
}

// SetAttributes устанавливает атрибуты модели из map
// func (u *User) SetAttributes(attrs map[string]interface{}) {
// 	for key, value := range attrs {
// 		switch key {
// 		case "id":
// 			if id, ok := value.(int64); ok {
// 				u.ID = int(id)
// 			}
// 		case "name":
// 			if name, ok := value.(string); ok {
// 				u.Name = name
// 			}
// 		case "email":
// 			if email, ok := value.(string); ok {
// 				u.Email = email
// 			}
// 		case "password":
// 			if password, ok := value.(string); ok {
// 				u.Password = password
// 			}
// 		case "role":
// 			if role, ok := value.(string); ok {
// 				u.Role = role
// 			}
// 		case "created_at":
// 			if createdAt, ok := value.(time.Time); ok {
// 				u.CreatedAt = createdAt
// 			}
// 		case "updated_at":
// 			if updatedAt, ok := value.(time.Time); ok {
// 				u.UpdatedAt = updatedAt
// 			}
// 		}
// 	}
// }

// UserRepository репозиторий для работы с пользователями
// type UserRepository struct {
// 	db *sql.DB
// }

// func NewUserRepository(db *sql.DB) *UserRepository {
// 	return &UserRepository{db: db}
// }

// RegisterRequest структура для запроса регистрации
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest структура для запроса входа
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}