package repository
import (
	"go_shopmarket/products/dto"
)

type Repository interface {
	CreateProduct(product dto.Products) error

}
