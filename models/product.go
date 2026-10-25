package models

import (
	"mime/multipart"
)

type Product struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	Name              string     `gorm:"type:varchar(100);not null" json:"name"`
	Slug              string     `gorm:"type:varchar(100);not null;unique" json:"slug"`
	Quantity          int        `gorm:"not null" json:"quantity"`
	Categories        []Category `gorm:"many2many:product_categories" json:"categories"`
	Price             float64    `json:"price"`
	Description       string     `gorm:"type:text" json:"description"`
	Sku               string     `gorm:"type:varchar(100);not null; unique;" json:"sku"`
	PublicationStatus string     `gorm:"type:varchar(20);not null;default:draft" json:"publication_status"`
	// Tags              []string   `json:"tags"`
	Image string `gorm:"type:text" json:"image"`
	Metadata
}

type CreateProductRequest struct {
	Name        string   `form:"name" binding:"required,min=3,max=100"`
	Slug        string   `form:"slug"`
	Quantity    int      `form:"quantity" binding:"required,gte=0"`
	Categories  []uint   `form:"categories" binding:"required,min=1,dive,gt=0"`
	Price       *float64 `form:"price" binding:"required,gte=0"`
	Description string   `form:"description"`
	Sku         string   `form:"sku"`
	// Tags        []string              `form:"tags"`
	Image *multipart.FileHeader `form:"image"`
}

type UpdateProductRequest struct {
	Name        *string  `form:"name" json:"name" binding:"omitempty,min=3,max=100"`
	Slug        *string  `form:"slug" json:"slug"`
	Quantity    *int     `form:"quantity" json:"quantity" binding:"omitempty,gte=0"`
	Price       *float64 `form:"price" json:"price" binding:"omitempty,gte=0"`
	Categories  *[]uint  `form:"categories" json:"categories"`
	Description *string  `form:"description" json:"description"`
	Sku         *string  `form:"sku" json:"sku"`
	// Tags        *[]string             `form:"tags"`
	Image *multipart.FileHeader `form:"image" json:"image"`
}

type ProductResponse struct {
	ID                uint               `json:"id"`
	Name              string             `json:"name"`
	Slug              string             `json:"slug"`
	Quantity          int                `json:"quantity"`
	Categories        []CategoryResponse `json:"categories"`
	Price             float64            `json:"price"`
	Description       string             `json:"description"`
	Sku               string             `json:"sku"`
	PublicationStatus string             `json:"publication_status"`
	// Tags              []string           `json:"tags"`
	Image string `json:"image"`
	Metadata
}

type UpdateQuantityPayload struct {
	Quantity int `json:"quantity"`
}

type UpdateProductPublicationStatusRequest struct {
	PublicationStatus string `json:"publication_status" binding:"required,oneof=active archived"`
}
