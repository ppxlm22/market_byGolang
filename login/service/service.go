package service

import (

	"time"

	"go_shopmarket/login/dto"
	"go_shopmarket/login/repository"
	"go_shopmarket/apperror"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

func NewService(r repository.Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) LoginUser(req dto.LoginRequest) (string, dto.UserResponse, error) {
	user, hashedPassword, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		log.Println("Error จาก Database:", err)
		return "", dto.UserResponse{}, apperror.NewNotFound("ชื่อผู้ใช้ไม่พบ")
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		log.Println("Error จาก Bcrypt:", err)
		return "", dto.UserResponse{}, apperror.NewBadRequest("รหัสผ่านไม่ถูกต้อง")
	}
	secretKey := os.Getenv("JWT_SECRET")
	jwtSecret := []byte(secretKey)
	claims := jwt.MapClaims{
		"user_id":  user.ID,	
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", dto.UserResponse{}, apperror.NewInternalServerError("เกิดข้อผิดพลาดในการสร้าง token")
	}
	
	return t, user, nil
}