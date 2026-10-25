package models

import "time"

type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	Slug      string    `json:"slug" gorm:"type:varchar(100);not null; unique"`
	Status    string    `json:"status" gorm:"type:varchar(20);not null; default:archived"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" gorm:"type:varchar(100);not null"`
	Slug string `json:"slug" gorm:"type:varchar(100);not null"`
}
