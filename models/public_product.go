package models

type PublicProduct struct {
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string     `gorm:"type:varchar(100);not null;unique" json:"slug"`
	Quantity    int        `gorm:"not null" json:"quantity"`
	Categories  []Category `gorm:"many2many:product_categories" json:"categories"`
	Price       float64    `json:"price"`
	Description string     `gorm:"type:text" json:"description"`
	Image       string     `gorm:"type:text" json:"image"`
}
