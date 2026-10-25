package repositories

import (
	"myapi/db"
	"myapi/models"
)

type PublicCategoriesRepository struct{}

func (r *PublicCategoriesRepository) FindAll() ([]models.PublicCategories, error) {
	var categories []models.Category

	err := db.DB.
		Where("status = ?", "active").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	publicCategories := make([]models.PublicCategories, 0, len(categories))

	for _, category := range categories {
		publicCategories = append(publicCategories, models.PublicCategories{
			Name: category.Name,
			Slug: category.Slug,
		})
	}

	return publicCategories, nil
}

func (r *PublicCategoriesRepository) FindBySlug(slug string) (*models.PublicCategories, error) {
	var categories models.Category

	err := db.DB.
		Where("slug = ? AND  status = ?", slug, "active").
		First(&categories).Error

	if err != nil {
		return nil, err
	}

	return &models.PublicCategories{
		Name: categories.Name,
		Slug: categories.Slug,
	}, err
}
