package models

import "time"

type Product struct {
	ID          uint       `gorm:"primaryKey"`
	Name        string     `gorm:"type:varchar(100);not null"`
	Quantity    int        `gorm:"not null"`
	Categories  []Category `gorm:"many2many:product_categories;"`
	Price       string     `gorm:"type:varchar(100);not null"`
	Description string     `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
