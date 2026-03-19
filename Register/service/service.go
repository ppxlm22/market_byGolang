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

func (s *service) RegisterUser(req dto.RegisterRequest) error {

	if req.Password == "" {
		return errors.New("กรุณากรอกรหัสผ่าน")
	}
	isDuplicate, err := s.repo.CheckUserExists(req.Username, req.Email)
	if isDuplicate {
		return errors.New("Username หรือ Email นี้มีผู้ใช้งานแล้ว")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return errors.New("เกิดข้อผิดพลาดในการเข้ารหัสรหัสผ่าน")
	}
	req.Password = string(hashedPassword)

	err = s.repo.Register(req)
	if err != nil {
		return err
	}
	return nil
}