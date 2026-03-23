package repository
import (
	"go_shopmarket/products/dto"
)

type Repository interface {
	CreateProduct(product dto.Products) error
	GetAllProducts() ([]dto.Products, error)
	GetProductByID(id int) (dto.Products, error)
	UpdateProduct(id int, product dto.Products) error
	DeleteProduct(id int) error
	GetCategoryByID(id int) (string, error)
	GetAllCategories() ([]dto.Category, error)
	DeductProductStock(id int, quantity int) error

}
