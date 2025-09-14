package services

import (
	"github.com/mymindmap/api/internal/http/requests/log"
	"github.com/mymindmap/api/models"
	"github.com/mymindmap/api/repository"
)

type LogService interface {
	List(limit, offset int) ([]models.Log, error)
	Get(id int) (*models.Log, error)
	Create(req log.CreateLogRequest) error
	Update(id int, req log.UpdateLogRequest) error
	Delete(id int) error
}

type logService struct {
	repo *repository.LogRepository
}

func NewLogService(repo *repository.LogRepository) LogService {
	return &logService{repo: repo}
}

func (s *logService) List(limit, offset int) ([]models.Log, error) {
	return s.repo.List(limit, offset)
}

func (s *logService) Get(id int) (*models.Log, error) {
	return s.repo.Get(id)
}

func (s *logService) Create(req log.CreateLogRequest) error {
	item := &models.Log{
		// TODO: маппинг req -> модель
		// Example fields based on your Log model:
		// Title:   req.Title,
		// Content: req.Content,
		// UserId:  req.UserId,
	}
	return s.repo.Create(item)
}

func (s *logService) Update(id int, req log.UpdateLogRequest) error {
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
	
	return s.repo.Update(id, updates)
}

func (s *logService) Delete(id int) error {
	return s.repo.Delete(id)
}