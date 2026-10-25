package services

import (
	"myapi/models"
	"myapi/repositories"
)

func NewPublicProductService() *PublicProductService {
	return &PublicProductService{
		repo: &repositories.PublicProductRepository{},
	}
}

type PublicProductService struct {
	repo *repositories.PublicProductRepository
}

func (s *PublicProductService) GetAll() ([]models.PublicProduct, error) {
	return s.repo.FindAll()
}

func (s *PublicProductService) GetBySlug(slug string) (*models.PublicProduct, error) {
	return s.repo.FindBySlug(slug)
}
