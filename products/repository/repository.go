package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"go_shopmarket/products/dto"
)
var ErrDBQuery = errors.New("database query error")
var ErrProductNotFound = errors.New("not found product in database")
var ErrCategoryNotFound = errors.New("not found category in database")

type repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r *repository) CreateProduct(product dto.Products) error {

	query := `INSERT INTO public.products (name, price, stock, category_id ) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBQuery, err)
	}
	return nil
}
func (r *repository) GetAllProducts() ([]dto.Products, error) {
	products := []dto.Products{}
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, p.update_at, COALESCE(c.name, 'ไม่ระบุ') AS category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id ASC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBQuery, err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var product dto.Products
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.UpdatedAt, &product.CategoryName)
		if err != nil {
			return nil, fmt.Errorf("fail to scan product row: %w", err)
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return products, nil
}

func (r *repository) GetProductByID(id int) (dto.Products, error) {
	var product dto.Products
	query := `SELECT id, name, price, stock, category_id FROM products WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.Products{}, ErrProductNotFound
		}
		return dto.Products{}, fmt.Errorf("%w: %v", ErrDBQuery, err)
	}
	return product, nil
}

func (r *repository) UpdateProduct(id int, product dto.Products) error {
	query := `UPDATE public.products SET name = $1, price = $2,stock = $3 ,category_id = $4 WHERE id = $5`
	_, err := r.DB.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBQuery, err)
	}
	return nil
}
func (r *repository) UpdateStock(id int, quantity int) error {
	query := `UPDATE public.products 
        SET stock = stock - $1, 
        update_at = CURRENT_TIMESTAMP 
        WHERE id = $2 AND stock >= $1`
	result, err := r.DB.Exec(query, quantity, id)
	if err != nil {
		return  fmt.Errorf("%w: %v", ErrDBQuery, err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (r *repository) DeleteProduct(id int) error {
	query := `DELETE FROM public.products WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return ErrProductNotFound
	}
	return nil
}
func (r *repository) GetCategoryByID(id int) (string, error) {
	var categoryName string
	query := `SELECT name FROM categories WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&categoryName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrCategoryNotFound
		}
		return "", fmt.Errorf("%w: %v", ErrDBQuery, err)
	}
	return categoryName, nil
}
func (r *repository) GetAllCategories() ([]dto.Category, error) {
	categories := []dto.Category{}
	query := `SELECT id, name FROM categories ORDER BY id ASC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBQuery, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cat dto.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, fmt.Errorf("fail to scan category row: %w", err)
		}
		categories = append(categories, cat)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return categories, nil
}
