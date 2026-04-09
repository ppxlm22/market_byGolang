package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"go_shopmarket/register/dto"
	"go_shopmarket/register/repository"
)

func NewService(r repository.Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) RegisterUser(req dto.RegisterRequest) (*dto.RegisterRespone, error) {

	if req.Password == "" {
		return nil, errors.New("กรุณากรอกรหัสผ่าน")
	}
	isDuplicate, err := s.repo.CheckUserExists(req.Username, req.Email)
	if isDuplicate {
		return nil, errors.New("Username หรือ Email นี้มีผู้ใช้งานแล้ว")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, errors.New("เกิดข้อผิดพลาดในการเข้ารหัสรหัสผ่าน")
	}
	req.Password = string(hashedPassword)

	user, err := s.repo.Register(req)
	if err != nil {
		return nil, err
	}
	return &dto.RegisterRespone{
		Username: user.Username,
		Email: user.Email,
		Message: "success",
	}, nil
}