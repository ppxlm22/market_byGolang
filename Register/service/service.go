package service

import (
	"go_shopmarket/apperror"
	"golang.org/x/crypto/bcrypt"
	"go_shopmarket/register/dto"
	"go_shopmarket/register/repository"
)

func NewService(r repository.Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) RegisterUser(req dto.RegisterRequest) (*dto.RegisterResponse, error) {

	isDuplicate, err := s.repo.CheckUserExists(req.Username, req.Email)
	if err != nil {
		return nil, apperror.NewInternal("เกิดข้อผิดพลาดในการตรวจสอบข้อมูลผู้ใช้")
	}	
	if isDuplicate {
		return nil, apperror.NewConflict("Username หรือ Email นี้มีผู้ใช้งานแล้ว")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, apperror.NewInternal("เกิดข้อผิดพลาดในการเข้ารหัสรหัสผ่าน")
	}
	req.Password = string(hashedPassword)

	user, err := s.repo.Register(req)
	if err != nil {
		return nil, err
	}
	return &dto.RegisterResponse{
		Username: user.Username,
		Email: user.Email,
		Message: "success",
	}, nil
}