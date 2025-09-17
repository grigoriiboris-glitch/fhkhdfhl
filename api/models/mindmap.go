package models

import (
	"time"
	"github.com/mymindmap/api/pkg/relations"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MindMap struct {
	relations.BaseModel
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Data      string    `json:"data" db:"data"`
	UserID    int       `json:"user_id" db:"user_id"`
	IsPublic  bool      `json:"is_public" db:"is_public"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	User   *User
	db *pgxpool.Pool
}
func (u *MindMap) GetID() interface{} {
    return u.ID
}
func (m *MindMap) GetTableName() string { return "mindmaps" }
func (m *MindMap) GetConnection() *pgxpool.Pool { return m.db }

func (m *MindMap) SetAttributes(attrs map[string]interface{}) {
	for key, value := range attrs {
		switch key {
		case "id":
			if id, ok := value.(int64); ok {
				m.ID = int(id)
			}
		case "title":
			if title, ok := value.(string); ok {
				m.Data = title
			}
		case "data":
			if data, ok := value.(string); ok {
				m.Data = data
			}
		case "user_id":
			if userID, ok := value.(int64); ok {
				m.UserID = int(userID)
			}
		case "created_at":
			if createdAt, ok := value.(time.Time); ok {
				m.CreatedAt = createdAt
			}
		case "updated_at":
			if updatedAt, ok := value.(time.Time); ok {
				m.UpdatedAt = updatedAt
			}
		}
	}
}

type CreateMindMapRequest struct {
	Title    string `json:"title" validate:"required"`
	Data     string `json:"data" validate:"required"`
	IsPublic bool   `json:"is_public"`
}

type UpdateMindMapRequest struct {
	Title    string `json:"title" validate:"required"`
	Data     string `json:"data" validate:"required"`
	IsPublic bool   `json:"is_public"`
} 