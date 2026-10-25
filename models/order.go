package models

import "time"

type Order struct {
	ID         uint        `json:"id" gorm:"primaryKey"`
	UserID     uint        `json:"-" gorm:"not null"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status" gorm:"not null;default:pending"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"-"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderItemPayload struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type OrderPayload struct {
	Items []OrderItemPayload `json:"items"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type OrderProductResponse struct {
	ID         uint               `json:"id"`
	Name       string             `json:"name"`
	Categories []CategoryResponse `json:"categories"`
}

type OrderItemResponse struct {
	ID       uint                 `json:"id"`
	Product  OrderProductResponse `json:"product"`
	Quantity int                  `json:"quantity"`
	Price    float64              `json:"price"`
}

type OrderResponse struct {
	ID         uint                `json:"id"`
	Items      []OrderItemResponse `json:"items"`
	TotalPrice float64             `json:"total_price"`
	Status     string              `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

func (o *Order) ToResponse() OrderResponse {
	items := make([]OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		categories := make([]CategoryResponse, len(item.Product.Categories))
		for j, cat := range item.Product.Categories {
			categories[j] = CategoryResponse{ID: cat.ID, Name: cat.Name}
		}

		items[i] = OrderItemResponse{
			ID: item.ID,
			Product: OrderProductResponse{
				ID:         item.Product.ID,
				Name:       item.Product.Name,
				Categories: categories,
			},
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	return OrderResponse{
		ID:         o.ID,
		Items:      items,
		TotalPrice: o.TotalPrice,
		Status:     o.Status,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
	}
}
