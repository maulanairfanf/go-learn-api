package services

import (
	"context"
	"errors"
	"fmt"
	"myapi/db"
	"myapi/models"
	"myapi/repositories"
	"myapi/storage"
	"myapi/utils"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func validateCategoryIDs(requested []uint, categories []models.Category) error {
	seen := make(map[uint]struct{}, len(requested))
	duplicates := make([]uint, 0)
	duplicateSeen := make(map[uint]struct{})
	for _, id := range requested {
		if _, exists := seen[id]; exists {
			if _, alreadyAdded := duplicateSeen[id]; !alreadyAdded {
				duplicates = append(duplicates, id)
				duplicateSeen[id] = struct{}{}
			}
			continue
		}
		seen[id] = struct{}{}
	}
	if len(duplicates) > 0 {
		return &CategoryValidationError{
			Field:        "categories",
			Code:         "DUPLICATE_CATEGORY_IDS",
			DuplicateIDs: duplicates,
		}
	}

	found := make(map[uint]struct{}, len(categories))
	for _, category := range categories {
		found[category.ID] = struct{}{}
	}
	invalid := make([]uint, 0)
	for _, id := range requested {
		if _, exists := found[id]; !exists {
			invalid = append(invalid, id)
		}
	}
	if len(invalid) > 0 {
		return &CategoryValidationError{
			Field:      "categories",
			Code:       "INVALID_CATEGORY_IDS",
			InvalidIDs: invalid,
		}
	}
	return nil
}

type ProductService struct {
	repo    *repositories.ProductRepository
	storage storage.ImageStorage
}

func NewProductService(imageStorage storage.ImageStorage) *ProductService {
	return &ProductService{
		repo:    &repositories.ProductRepository{},
		storage: imageStorage,
	}
}

func (s *ProductService) GetAll(pagination utils.PaginationRequest, search string, sort string, order string) ([]models.Product, int64, error) {
	return s.repo.FindAll(pagination, search, sort, order)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) Create(ctx context.Context, req models.CreateProductRequest) (*models.Product, error) {
	if req.Price == nil || req.Image == nil {
		return nil, ErrInvalidInput
	}

	for _, categoryID := range req.Categories {
		if categoryID == 0 {
			return nil, ErrInvalidInput
		}
	}

	categories, err := s.repo.FindCategoriesByIDs(req.Categories)
	if err != nil {
		return nil, ErrInvalidInput
	}

	if err := validateCategoryIDs(req.Categories, categories); err != nil {
		return nil, err
	}

	file, err := req.Image.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	contentType := req.Image.Header.Get("Content-Type")
	imageURL, err := s.storage.Upload(ctx, req.Image.Filename, contentType, req.Image.Size, file)
	if err != nil {
		return nil, err
	}

	product := models.Product{
		Name:        req.Name,
		Slug:        req.Slug,
		Quantity:    req.Quantity,
		Categories:  categories,
		Price:       *req.Price,
		Description: req.Description,
		Sku:         req.Sku,
		Image:       imageURL,
		// Tags:        req.Tags,
	}

	if err := s.repo.Create(&product); err != nil {
		return nil, mapProductConflict(err)
	}

	s.repo.PreloadCategories(&product)
	return &product, nil
}

func (s *ProductService) Update(ctx context.Context, id int, req models.UpdateProductRequest) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		product.Name = *req.Name
	}

	if req.Quantity != nil {
		product.Quantity = *req.Quantity
	}

	if req.Price != nil {
		product.Price = *req.Price
	}

	if req.Description != nil {
		product.Description = *req.Description
	}

	if req.Slug != nil {
		exists, err := s.repo.ExistsSlug(*req.Slug, product.ID)
		fmt.Println("exists", exists)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, &ProductConflictError{
				Field:   "slug",
				Message: "Product slug already exists",
			}
		}
		product.Slug = *req.Slug
	}

	if req.Sku != nil {
		exists, err := s.repo.ExistsSku(*req.Sku, product.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, &ProductConflictError{
				Field:   "sku",
				Message: "Product sku already exists",
			}
		}
		product.Sku = *req.Sku
	}

	if req.Categories != nil {
		reqCategories := *req.Categories
		if len(reqCategories) > 0 {
			categories, err := s.repo.FindCategoriesByIDs(reqCategories)
			if err != nil {
				return nil, ErrInvalidInput
			}
			if err := validateCategoryIDs(reqCategories, categories); err != nil {
				return nil, err
			}
			if err := s.repo.ReplaceCategories(product, categories); err != nil {
				return nil, err
			}
		}

	}

	var oldImageUrl string

	if req.Image != nil {
		file, err := req.Image.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		contentType := req.Image.Header.Get("Content-Type")
		imageUrl, err := s.storage.Upload(ctx, req.Image.Filename, contentType, req.Image.Size, file)
		if err != nil {
			return nil, err
		}

		oldImageUrl = product.Image
		product.Image = imageUrl
	}

	if err := s.repo.Save(product); err != nil {
		return nil, mapProductConflict(err)
	}

	if product.Image != "" {
		err = s.storage.Delete(ctx, oldImageUrl)
		if err != nil {
			return nil, err
		}
	}

	s.repo.PreloadCategories(product)
	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id int) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(product).
			Association("Categories").
			Clear(); err != nil {
			return err
		}
		if err = tx.Delete(product).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = s.storage.Delete(ctx, product.Image)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) UpdatePublicationStatus(ctx context.Context, id int, req models.UpdateProductPublicationStatusRequest) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.PublicationStatus != "active" &&
		req.PublicationStatus != "archived" {
		return nil, ErrInvalidStatus
	}

	product.PublicationStatus = req.PublicationStatus

	err = s.repo.Save(product)
	if err != nil {
		return nil, err
	}

	return product, err
}

func mapProductConflict(err error) error {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return err
	}

	switch pgErr.ConstraintName {
	case "uni_products_slug":
		return &ProductConflictError{
			Field:   "slug",
			Message: "Product slug already exists",
		}
	case "uni_products_sku":
		return &ProductConflictError{
			Field:   "sku",
			Message: "Product SKU already exists",
		}
	default:
		return ErrConflict
	}
}
