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
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock"`
	CategoryID   int     `json:"category_id"`
	UpdatedAt    *time.Time `json:"updated_at"`
	CategoryName string  `json:"category"`
}
