package services

import (
	"github.com/mymindmap/api/internal/http/requests/mindMap"
	"github.com/mymindmap/api/repository"
)

type MindMapService interface {
	List() ([]*repository.MindMap, error)
	Get(id int) (*repository.MindMap, error)
	Create(req MindMapCreateRequest) (*repository.MindMap, error)
	Update(id int, req MindMapUpdateRequest) (*repository.MindMap, error)
	Delete(id int) error
}

type mindMapService struct {
	repo repository.MindMapRepository
}

func NewMindMapService(repo repository.MindMapRepository) MindMapService {
	return &mindMapService{repo: repo}
}

func (s *mindMapService) List() ([]*repository.MindMap, error) {
	return s.repo.List()
}

func (s *mindMapService) Get(id int) (*repository.MindMap, error) {
	return s.repo.Get(id)
}

func (s *mindMapService) Create(req MindMapCreateRequest) (*repository.MindMap, error) {
	item := &repository.MindMap{
		// TODO: маппинг req -> модель
	}
	return s.repo.Create(item)
}

func (s *mindMapService) Update(id int, req MindMapUpdateRequest) (*repository.MindMap, error) {
	item, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	// TODO: обновить поля item из req
	return s.repo.Update(item)
}

func (s *mindMapService) Delete(id int) error {
	return s.repo.Delete(id)
}
