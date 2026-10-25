package repositories

import (
	"myapi/db"
	"myapi/models"
	"myapi/utils"
)

type ProductRepository struct{}

func (r *ProductRepository) FindAll(pagination utils.PaginationRequest, search string, sort string, order string) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := db.DB.Model(&models.Product{})

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if sort != "" {
		query = query.Order(sort + " " + order)
	}

	query = utils.ApplyPagination(query, pagination, &total)
	err := query.Preload("Categories").Find(&products).Error

	return products, total, err
}

func (r *ProductRepository) FindByID(id int) (*models.Product, error) {
	var product models.Product
	err := db.DB.Preload("Categories").First(&product, id).Error
	return &product, err
}

func (r *ProductRepository) FindCategoriesByIDs(ids []uint) ([]models.Category, error) {
	var categories []models.Category
	err := db.DB.Where("id IN ?", ids).Find(&categories).Error
	return categories, err
}

func (r *ProductRepository) Create(product *models.Product) error {
	return db.DB.Create(product).Error
}

func (r *ProductRepository) ReplaceCategories(product *models.Product, categories []models.Category) error {
	return db.DB.Model(product).Association("Categories").Replace(categories)
}

func (r *ProductRepository) Save(product *models.Product) error {
	return db.DB.Save(product).Error
}

func (r *ProductRepository) ClearCategories(product *models.Product) error {
	return db.DB.Model(product).Association("Categories").Clear()
}

func (r *ProductRepository) Delete(product *models.Product) error {
	return db.DB.Delete(product).Error
}

func (r *ProductRepository) PreloadCategories(product *models.Product) error {
	return db.DB.Preload("Categories").First(product, product.ID).Error
}

func (r *ProductRepository) ExistsSlug(slug string, excludeID uint) (bool, error) {
	var count int64

	err := db.DB.Model(&models.Product{}).
		Where("slug = ? AND id <> ?", slug, excludeID).
		Count(&count).Error

	return count > 0, err
}

func (r *ProductRepository) ExistsSku(sku string, excludeID uint) (bool, error) {
	var count int64

	err := db.DB.Model(&models.Product{}).
		Where("sku = ? AND id <> ?", sku, excludeID).
		Count(&count).Error

	return count > 0, err
}
