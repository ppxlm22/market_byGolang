package login

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	LoginUser(req LoginRequest) (string, error)
}
type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) LoginUser(req LoginRequest) (string, error) {
	userId, hashedPassword, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return "", errors.New("ชื่อผู้ใช้ หรือ รหัสไม่ถูกต้อง")
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		return "",errors.New("รหัสผ่านไม่ถูกต้อง")
	}
	jwtSecret := []byte("your_secret_key")
	claims := jwt.MapClaims{
		"user_id":  userId,
		"username": req.Username,
		"role":     "customer",
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errors.New("ไม่สามารถสร้าง Token ได้")
	}
	return t, nil
}