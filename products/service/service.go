package service

import (
	"go_shopmarket/products/dto"
	"go_shopmarket/products/repository"
	"go_shopmarket/apperror"
)


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
	product, err := s.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}
	if len(product) == 0 {
		return nil, apperror.NewNotFound("ไม่พบสินค้า")
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

func (s *service) UpdateProduct(id int, product dto.Products) error {
	_, err := s.repo.GetProductByID(id)
	if err != nil {
		return apperror.NewNotFound("ไม่พบสินค้าที่ต้องการอัพเดต")
	}
	return s.repo.UpdateProduct(id, product)
}

func (s *service) DeleteProduct(id int) error {
	_, err := s.repo.GetProductByID(id)
	if err != nil {
		return apperror.NewNotFound("ไม่พบสินค้าที่ต้องการลบ")
	}
	return s.repo.DeleteProduct(id)
}

func (s *service) GetCategoryByID(id int) (string, error) {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return "", err
	}
	return category, nil
}

func (s *service) GetAllCategories() ([]dto.Category, error) {
	categories, err := s.repo.GetAllCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *service) Checkout(req dto.CheckoutRequest) error {
	if len(req.Items) == 0 {	
		return apperror.NewBadRequest("ไม่มีสินค้าในคำสั่งซื้อ")
	}
	for _, item := range req.Items {
		err := s.repo.DeductProductStock(item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}
