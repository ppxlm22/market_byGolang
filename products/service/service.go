package service

import (
	"go_shopmarket/products/dto"
	"go_shopmarket/products/repository"
	"errors"
	"fmt"
)
var ErrProductNotFound = errors.New("not found product in database")
var ErrEmptyCart = errors.New("no items in the cart for checkout")

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateProduct(product dto.Products) error {
	return s.repo.CreateProduct(product)
}

func (s *service) GetAllProducts() ([]dto.Products, error) {
	products, err := s.repo.GetAllProducts()
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}
	return products, nil
}

func (s *service) GetProductByID(id int) (dto.Products, error) {
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		return dto.Products{}, fmt.Errorf("failed to get product id %d: %w", id, err)
	}
	return product, nil
}

func (s *service) UpdateProduct(id int, product dto.Products) error {
	_, err := s.repo.GetProductByID(id)
	if err != nil {
		return fmt.Errorf("failed to update product id %d: %w", id, err)
	}
	return s.repo.UpdateProduct(id, product)
}

func (s *service) DeleteProduct(id int) error {
	_, err := s.repo.GetProductByID(id)
	if err != nil {
		return fmt.Errorf("failed to delete product id %d: %w", id, err)
	}
	return s.repo.DeleteProduct(id)
}

func (s *service) GetCategoryByID(id int) (string, error) {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return "", fmt.Errorf("failed to get category id %d: %w", id, err)
	}
	return category, nil
}

func (s *service) GetAllCategories() ([]dto.Category, error) {
	categories, err := s.repo.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}
	return categories, nil
}

func (s *service) Checkout(req dto.CheckoutRequest) error {
	if len(req.Items) == 0 {	
		return ErrEmptyCart
	}
	for _, item := range req.Items {
		err := s.repo.UpdateStock(item.ProductID, item.Quantity)
		if err != nil {
			return fmt.Errorf("checkout failed on product_id %d: %w", item.ProductID, err)
		}
	}
	return nil
}
