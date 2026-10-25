package repositories

import (
	"myapi/db"
	"myapi/models"
)

type PublicProductRepository struct{}

func (r *PublicProductRepository) FindAll() ([]models.PublicProduct, error) {
	var products []models.Product

	err := db.DB.
		Where("publication_status = ?", "active").
		Preload("Categories").
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	publicProducts := make([]models.PublicProduct, 0, len(products))

	for _, product := range products {
		publicProducts = append(publicProducts, models.PublicProduct{
			Name:        product.Name,
			Slug:        product.Slug,
			Categories:  product.Categories,
			Price:       product.Price,
			Description: product.Description,
			Image:       product.Image,
		})
	}

	return publicProducts, nil
}

func (r *PublicProductRepository) FindBySlug(slug string) (*models.PublicProduct, error) {
	var product models.Product

	err := db.DB.
		Where("slug = ? AND publication_status = ?", slug, "active").
		Preload("Categories").
		First(&product).Error

	if err != nil {
		return nil, err
	}

	return &models.PublicProduct{
		Name:        product.Name,
		Slug:        product.Slug,
		Categories:  product.Categories,
		Price:       product.Price,
		Description: product.Description,
		Image:       product.Image,
	}, nil
}
