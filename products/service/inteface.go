package service

import (
	"go_shopmarket/products/dto"
)

type Service interface {
	CreateProduct(product dto.Products) error
}

