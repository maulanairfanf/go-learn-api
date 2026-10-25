package repositories

import (
	"myapi/db"
	"myapi/models"
	"myapi/utils"
)

type CategoryRepository struct{}

func (r *CategoryRepository) FindAll(pagination utils.PaginationRequest, search string, sort string, order string) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	query := db.DB.Model(&models.Category{})

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if sort != "" {
		query = query.Order(sort + " " + order)
	}

	query = utils.ApplyPagination(query, pagination, &total)
	err := query.Find(&categories).Error

	if !pagination.IsPaginated() {
		total = int64(len(categories))
	}

	return categories, total, err
}

func (r *CategoryRepository) FindByID(id int) (*models.Category, error) {
	var category models.Category
	err := db.DB.First(&category, id).Error
	return &category, err
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return db.DB.Create(category).Error
}

func (r *CategoryRepository) Save(category *models.Category) error {
	return db.DB.Save(category).Error
}

func (r *CategoryRepository) Delete(category *models.Category) error {
	return db.DB.Delete(category).Error
}
