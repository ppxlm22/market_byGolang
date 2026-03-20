package service 

import (
	"errors"
	"go_shopmarket/products/dto"
	"go_shopmarket/products/repository"
)

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateProduct(product dto.Products) error {
	if product.Name == "" || product.Price <= 0 || product.Stock < 0 || product.CategoryID == 0 {
		return errors.New("ข้อมูลสินค้าไม่ถูกต้อง")
	}
	return s.repo.CreateProduct(product)
}
func (s *service) GetAllProducts() ([]dto.Products, error) {
	product, err := s.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}
	if len(product) == 0 {
		return nil, errors.New("ไม่พบสินค้า")
	}
	return product, nil
}	
func (s *service) GetProductByID(id int) (dto.Products, error) {
	product, err := s.repo.GetProductByID(id)	
	if err != nil {
		return dto.Products{}, err
	}
	return product, nil
}

