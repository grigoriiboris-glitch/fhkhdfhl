package services

import (
	"github.com/mymindmap/api/internal/http/requests/user"
	"github.com/mymindmap/api/repository"
)

type UserService interface {
	List() ([]*repository.User, error)
	Get(id int) (*repository.User, error)
	Create(req UserCreateRequest) (*repository.User, error)
	Update(id int, req UserUpdateRequest) (*repository.User, error)
	Delete(id int) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) List() ([]*repository.User, error) {
	return s.repo.List()
}

func (s *userService) Get(id int) (*repository.User, error) {
	return s.repo.Get(id)
}

func (s *userService) Create(req UserCreateRequest) (*repository.User, error) {
	item := &repository.User{
		// TODO: маппинг req -> модель
	}
	return s.repo.Create(item)
}

func (s *userService) Update(id int, req UserUpdateRequest) (*repository.User, error) {
	item, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	// TODO: обновить поля item из req
	return s.repo.Update(item)
}

func (s *userService) Delete(id int) error {
	return s.repo.Delete(id)
}
