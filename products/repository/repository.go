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
func (r *repository) GetProductByID(id int) (dto.Products, error) {
	var product dto.Products
	query := `SELECT id, name, price, stock, category_id FROM products WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID)
	if err != nil {
		return dto.Products{}, err
	}
	return product, nil
}
func (r *repository) UpdateProduct(id int, product dto.Products) error {
	query := `UPDATE public.products SET name = $1, price = $2,stock = $3 ,category_id = $4 WHERE id = $5`
	_, err := database.DB.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, id)
	if err != nil {
		return err	
	}
	return nil
}