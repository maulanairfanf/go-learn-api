package models

import "time"

type Category struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

