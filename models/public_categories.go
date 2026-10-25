package models

type PublicCategories struct {
	Name string `json:"name" gorm:"type:varchar(100);not null"`
	Slug string `json:"slug" gorm:"type:varchar(100);not null; unique"`
}
