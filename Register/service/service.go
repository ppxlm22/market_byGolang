package service

import (
	"golang.org/x/crypto/bcrypt"
	"go_shopmarket/register/dto"
	"go_shopmarket/register/repository"
	"errors"
	"fmt"
)
var ErrUserDuplicate = errors.New("username or email already exists")

func NewService(r repository.Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) RegisterUser(req dto.RegisterRequest) (*dto.RegisterResponse, error) {

	isDuplicate, err := s.repo.CheckUserExists(req.Username, req.Email)
	if err != nil {
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการตรวจสอบข้อมูลผู้ใช้: %w", err)
	}	
	if isDuplicate {
		return nil, ErrUserDuplicate
	}
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil,fmt.Errorf("เกิดข้อผิดพลาดในการเข้ารหัสผ่าน: %w", err)
	}
	req.Password = hashedPassword

	user, err := s.repo.Register(req)
	if err != nil {
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการลงทะเบียนผู้ใช้: %w", err)
	}
	return &dto.RegisterResponse{
		Username: user.Username,
		Email: user.Email,
		Message: "success",
	}, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}