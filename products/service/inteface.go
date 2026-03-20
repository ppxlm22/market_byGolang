package service

import (
	"go_shopmarket/products/dto"
)

type Service interface {
	CreateProduct(product dto.Products) error
	GetAllProducts() ([]dto.Products, error)
	GetProductByID(id int) (dto.Products, error)
}
