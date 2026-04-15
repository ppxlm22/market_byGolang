package dto

import(
	"time"
)

type Category struct {
	ID int `json:"id"`
	Name string `json:"name"`
}
type Products struct {
	ID           int     `json:"id"`
	Name         string  `json:"name" validate:"required"`
	Price        float64 `json:"price" validate:"required,gt=0"`
	Stock        int     `json:"stock" validate:"required,gte=0"`
	CategoryID   int     `json:"category_id" validate:"required"`
	UpdatedAt    *time.Time `json:"update_at"`
	CategoryName string  `json:"category"`
}

type checkoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity" validate:"required,gt=0"`
}
type CheckoutRequest struct {
	Items []checkoutItem `json:"items"`
}

