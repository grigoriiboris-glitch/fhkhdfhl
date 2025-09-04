package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mymindmap/api/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MindMapRepository struct {
	db *pgxpool.Pool
}

func NewMindMapRepository(db *pgxpool.Pool) *MindMapRepository {
	return &MindMapRepository{db: db}
}

func (r *MindMapRepository) Create(ctx context.Context, mindMap *models.MindMap) error {
	query := `
		INSERT INTO mindmaps (title, data, user_id, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		mindMap.Title,
		mindMap.Data,
		mindMap.UserID,
		mindMap.IsPublic,
		now,
		now,
	).Scan(&mindMap.ID)

	if err != nil {
		return fmt.Errorf("error creating mindmap: %w", err)
	}

	mindMap.CreatedAt = now
	mindMap.UpdatedAt = now
	return nil
}

func (r *MindMapRepository) GetByID(ctx context.Context, id int) (*models.MindMap, error) {
	query := `
		SELECT id, title, data, user_id, is_public, created_at, updated_at
		FROM mindmaps
		WHERE id = $1`

	mindMap := &models.MindMap{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&mindMap.ID,
		&mindMap.Title,
		&mindMap.Data,
		&mindMap.UserID,
		&mindMap.IsPublic,
		&mindMap.CreatedAt,
		&mindMap.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting mindmap: %w", err)
	}

	return mindMap, nil
}

func (r *MindMapRepository) GetByUserID(ctx context.Context, userID int) ([]*models.MindMap, error) {
	query := `
		SELECT id, title, data, user_id, is_public, created_at, updated_at
		FROM mindmaps
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get mindmaps by user: %w", err)
	}
	defer rows.Close()

	var mindMaps []*models.MindMap

	for rows.Next() {
		mindMap := new(models.MindMap) // то же самое, что &models.MindMap{}
		if err := rows.Scan(
			&mindMap.ID,
			&mindMap.Title,
			&mindMap.Data,
			&mindMap.UserID,
			&mindMap.IsPublic,
			&mindMap.CreatedAt,
			&mindMap.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan mindmap: %w", err)
		}
		mindMaps = append(mindMaps, mindMap)
	}

	// очень важно: проверяем ошибки после итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return mindMaps, nil
}


func (r *MindMapRepository) GetPublic(ctx context.Context) ([]*models.MindMap, error) {
	query := `
		SELECT id, title, data, user_id, is_public, created_at, updated_at
		FROM mindmaps
		WHERE is_public = true
		ORDER BY updated_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting public mindmaps: %w", err)
	}
	defer rows.Close()

	var mindMaps []*models.MindMap
	for rows.Next() {
		mindMap := &models.MindMap{}
		err := rows.Scan(
			&mindMap.ID,
			&mindMap.Title,
			&mindMap.Data,
			&mindMap.UserID,
			&mindMap.IsPublic,
			&mindMap.CreatedAt,
			&mindMap.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning mindmap: %w", err)
		}
		mindMaps = append(mindMaps, mindMap)
	}

	return mindMaps, nil
}

func (r *MindMapRepository) Update(ctx context.Context, mindMap *models.MindMap) error {
	query := `
		UPDATE mindmaps
		SET title = $1, data = $2, is_public = $3, updated_at = $4
		WHERE id = $5 AND user_id = $6`

	result, err := r.db.Exec(ctx, query,
		mindMap.Title,
		mindMap.Data,
		mindMap.IsPublic,
		time.Now(),
		mindMap.ID,
		mindMap.UserID,
	)

	if err != nil {
		return fmt.Errorf("error updating mindmap: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("mindmap not found or access denied")
	}

	return nil
}

func (r *MindMapRepository) Delete(ctx context.Context, id, userID int) error {
	query := `DELETE FROM mindmaps WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("error deleting mindmap: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("mindmap not found or access denied")
	}

	return nil
}

func (r *MindMapRepository) DeleteByAdmin(ctx context.Context, id int) error {
	query := `DELETE FROM mindmaps WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting mindmap: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("mindmap not found")
	}

	return nil
} 

// Compatibility methods expected by internal/handlers/mindmaps.go
// They delegate to the existing repository methods defined above.

// GetMindMapsByUser returns mindmaps for a given user.
func (r *MindMapRepository) GetMindMapsByUser(ctx context.Context, userID int) ([]*models.MindMap, error) {
    return r.GetByUserID(ctx, userID)
}

// GetMindMapByID returns a single mindmap by ID.
func (r *MindMapRepository) GetMindMapByID(ctx context.Context, id int) (*models.MindMap, error) {
    return r.GetByID(ctx, id)
}

// CreateMindMap creates a new mindmap.
func (r *MindMapRepository) CreateMindMap(ctx context.Context, m *models.MindMap) error {
    return r.Create(ctx, m)
}

// UpdateMindMap updates an existing mindmap.
func (r *MindMapRepository) UpdateMindMap(ctx context.Context, m *models.MindMap) error {
    return r.Update(ctx, m)
}

// DeleteMindMap deletes a mindmap by ID (and optionally checks user elsewhere).
func (r *MindMapRepository) DeleteMindMap(ctx context.Context, id int) error {
    // For generic deletion without user constraint, reuse admin delete
    return r.DeleteByAdmin(ctx, id)
}