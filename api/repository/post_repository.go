package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/betterstack-community/go-blog/models"
)

type PostRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostRepository(dbpool *pgxpool.Pool) *PostRepository {
	return &PostRepository{dbpool}
}

func (pr *PostRepository) CreatePost(ctx context.Context, post *models.Post) error {
	query := `
		INSERT INTO posts (title, content, user_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`

	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now

	return pr.dbpool.QueryRow(ctx, query, 
		post.Title, 
		post.Content, 
		post.UserID,
		post.CreatedAt,
		post.UpdatedAt,
	).Scan(&post.ID)
}

func (pr *PostRepository) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	rows, err := pr.dbpool.Query(
		ctx,
		"SELECT id, title, content, user_id, created_at, updated_at FROM posts ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID, 
			&post.Title, 
			&post.Content,
			&post.UserID,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (pr *PostRepository) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	var post models.Post
	query := `
		SELECT id, title, content, user_id, created_at, updated_at 
		FROM posts 
		WHERE id = $1`
	
	err := pr.dbpool.QueryRow(ctx, query, id).Scan(
		&post.ID, 
		&post.Title, 
		&post.Content,
		&post.UserID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &post, nil
}

func (pr *PostRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	query := `
		UPDATE posts 
		SET title = $1, content = $2, updated_at = $3 
		WHERE id = $4`
	
	post.UpdatedAt = time.Now()
	
	_, err := pr.dbpool.Exec(
		ctx,
		query,
		post.Title,
		post.Content,
		post.UpdatedAt,
		post.ID,
	)
	return err
}

func (pr *PostRepository) DeletePost(ctx context.Context, id int) error {
	query := "DELETE FROM posts WHERE id = $1"
	_, err := pr.dbpool.Exec(ctx, query, id)
	return err
}

// GetPostsByUser получает посты конкретного пользователя
func (pr *PostRepository) GetPostsByUser(ctx context.Context, userID int) ([]*models.Post, error) {
	rows, err := pr.dbpool.Query(
		ctx,
		"SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID, 
			&post.Title, 
			&post.Content,
			&post.UserID,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}
