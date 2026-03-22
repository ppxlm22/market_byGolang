package repository

import (
	"go_shopmarket/database"
	"go_shopmarket/products/dto"
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
	query := `
		SELECT 
			p.id, 
			p.name, 
			p.price, 
			p.stock, 
			p.category_id, 
			p.update_at,
			COALESCE(c.name, 'ไม่ระบุ') AS category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
        ORDER BY p.id ASC`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product dto.Products
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.UpdatedAt, &product.CategoryName)
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
func (r *repository) DeleteProduct(id int) error {
	query := `DELETE FROM public.products WHERE id = $1`
	_, err := database.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
func (r *repository) GetCategoryByID(id int) (string, error) {
	var categoryName string
	query := `SELECT name FROM categories WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(&categoryName)
	if err != nil {
		return "", err
	}
	return categoryName, nil
}
func (r *repository) GetAllCategories() ([]dto.Category, error) {
	var categories []dto.Category
	query := `SELECT id, name FROM categories ORDER BY id ASC`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cat dto.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}
