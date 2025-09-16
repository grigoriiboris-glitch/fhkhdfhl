package log_service

import (
	"context"
	"github.com/mymindmap/api/internal/http/requests/log"
	"github.com/mymindmap/api/models"
	"github.com/mymindmap/api/repository"
)

type LogService interface {
	List(ctx context.Context, limit, offset int) ([]models.Log, error)
	Get(ctx context.Context, id int) (*models.Log, error)
	Create(ctx context.Context, req log.CreateLogRequest) error
	Update(ctx context.Context, id int, req log.UpdateLogRequest) error
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}

type logService struct {
	repo *repository.LogRepository
}

func NewLogService(repo *repository.LogRepository) LogService {
	return &logService{repo: repo}
}

func (s *logService) List(ctx context.Context, limit, offset int) ([]models.Log, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *logService) Get(ctx context.Context, id int) (*models.Log, error) {
	return s.repo.Get(ctx, id)
}

func (s *logService) Create(ctx context.Context, req log.CreateLogRequest) error {
	item := &models.Log{
		// TODO: маппинг req -> модель
		// Example fields based on your Log model:
		// Title:   req.Title,
		// Content: req.Content,
		// UserId:  req.UserId,
	}
	return s.repo.Create(ctx, item)
}

func (s *logService) Update(ctx context.Context, id int, req log.UpdateLogRequest) error {
	// Create update map based on request
	updates := make(map[string]interface{})
	
	// TODO: Add field mappings from req to updates map
	// Example:
	// if req.Title != "" {
	//     updates["title"] = req.Title
	// }
	// if req.Content != "" {
	//     updates["content"] = req.Content
	// }
	// if req.UserId != 0 {
	//     updates["user_id"] = req.UserId
	// }
	
	return s.repo.Update(ctx, id, updates)
}

func (s *logService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *logService) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}