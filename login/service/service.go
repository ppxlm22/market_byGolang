package service

import (
	"errors"
	"time"

	"go_shopmarket/login/dto"
	"go_shopmarket/login/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func NewService(r repository.Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) LoginUser(req dto.LoginRequest) (string, error) {
	userID,hashedPassword, userRole, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		log.Println("Error จาก Database:", err)
		return "", errors.New("ชื่อผู้ใช้ หรือ รหัสไม่ถูกต้อง")
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		log.Println("Error จาก Bcrypt:", err)
		return "", errors.New("รหัสผ่านไม่ถูกต้อง")
	}
	jwtSecret := []byte("your_secret_key")
	claims := jwt.MapClaims{
		"user_id":  userID,	
		"username": req.Username,
		"role":     userRole,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errors.New("ไม่สามารถสร้าง Token ได้")
	}
	
	return t, nil
}