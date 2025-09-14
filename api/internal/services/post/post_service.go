package services

import (
	"github.com/mymindmap/api/internal/http/requests/post"
	"github.com/mymindmap/api/repository"
)

type PostService interface {
	List() ([]*repository.Post, error)
	Get(id int) (*repository.Post, error)
	Create(req PostCreateRequest) (*repository.Post, error)
	Update(id int, req PostUpdateRequest) (*repository.Post, error)
	Delete(id int) error
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) List() ([]*repository.Post, error) {
	return s.repo.List()
}

func (s *postService) Get(id int) (*repository.Post, error) {
	return s.repo.Get(id)
}

func (s *postService) Create(req PostCreateRequest) (*repository.Post, error) {
	item := &repository.Post{
		// TODO: маппинг req -> модель
	}
	return s.repo.Create(item)
}

func (s *postService) Update(id int, req PostUpdateRequest) (*repository.Post, error) {
	item, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	// TODO: обновить поля item из req
	return s.repo.Update(item)
}

func (s *postService) Delete(id int) error {
	return s.repo.Delete(id)
}
