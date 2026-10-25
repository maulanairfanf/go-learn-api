package services

import (
	"myapi/models"
	"myapi/repositories"
	"myapi/utils"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		repo: &repositories.CategoryRepository{},
	}
}

func (s *CategoryService) GetAll(pagination utils.PaginationRequest, search string, sort string, order string) ([]models.Category, int64, error) {
	return s.repo.FindAll(pagination, search, sort, order)
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) Create(req models.CreateCategoryRequest) (*models.Category, error) {
	category := models.Category{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := s.repo.Create(&category); err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *CategoryService) Update(id int, req models.CreateCategoryRequest) (*models.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Slug = req.Slug

	if err := s.repo.Save(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) Delete(id int) (*models.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Delete(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) ChangeStatus(id int) (*models.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if category.Status == "active" {
		category.Status = "inactive"
	} else {
		category.Status = "active"
	}

	if err := s.repo.Save(category); err != nil {
		return nil, err
	}

	return category, nil
}
