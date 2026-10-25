package services

import (
	"log/slog"
	"myapi/models"
	"myapi/repositories"
	"myapi/utils"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: &repositories.UserRepository{},
	}
}

func (s *UserService) GetAll(pagination utils.PaginationRequest, search string, sort string, order string) ([]models.User, int64, error) {
	return s.repo.FindAll(pagination, search, sort, order)
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) Delete(id int) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	err = s.repo.Delete(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Update(id int, req models.UpdateUserRequest) (*models.User, error) {
	slog.Info("request update", "username", req.Username)
	_, err := s.repo.FindByUsername(req.Username)
	if err == nil {
		return nil, ErrConflict
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	user.Username = req.Username

	err = s.repo.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
