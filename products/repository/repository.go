package repository
import (
	"go_shopmarket/products/dto"
	"go_shopmarket/database"
)

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreateProduct(product dto.Products) error {

	query := `INSERT INTO public.products (name, price, stock, category_id ) VALUES ($1, $2, $3, $4)`
	_, err := database.DB.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID)
	if err != nil {
		return err
	}
	return nil
}
func (r *repository) GetAllProducts() ([]dto.Products, error) {
	var products []dto.Products
	query := `SELECT id, name, price, stock, category_id FROM products`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product dto.Products
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
