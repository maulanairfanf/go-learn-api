package services

import (
	"myapi/models"
	"myapi/repositories"
)

func NewPublicCategoriesService() *PublicCategoriesService {
	return &PublicCategoriesService{
		repo: &repositories.PublicCategoriesRepository{},
	}
}

type PublicCategoriesService struct {
	repo *repositories.PublicCategoriesRepository
}

func (s *PublicCategoriesService) GetAll() ([]models.PublicCategories, error) {
	return s.repo.FindAll()
}

func (s *PublicCategoriesService) GetBySlug(slug string) (*models.PublicCategories, error) {
	return s.repo.FindBySlug(slug)
}
