package register

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)
type Service interface {
	RegisterUser(req registerDB) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) RegisterUser(req registerDB) error {
	if req.Password == "" {
		return errors.New("กรุณากรอกรหัสผ่าน")
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