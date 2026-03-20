package repository
import (
	"go_shopmarket/products/dto"
)

type Repository interface {
	CreateProduct(product dto.Products) error
	GetAllProducts() ([]dto.Products, error)
	GetProductByID(id int) (dto.Products, error)

}
