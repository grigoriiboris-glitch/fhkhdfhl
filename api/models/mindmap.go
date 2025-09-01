package models

import (
	"time"
)

type MindMap struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Data      string    `json:"data" db:"data"`
	UserID    int       `json:"user_id" db:"user_id"`
	IsPublic  bool      `json:"is_public" db:"is_public"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
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