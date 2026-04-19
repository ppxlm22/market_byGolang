package service

import (
	"errors"
	"fmt"
	"time"

	"go_shopmarket/login/dto"
	"go_shopmarket/login/repository"

	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)
var ErrInvalidCredentials = errors.New("Username or password is incorrect")

func NewService(r repository.Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) LoginUser(req dto.LoginRequest) (string, dto.UserResponse, error) {
	user, hashedPassword, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", dto.UserResponse{}, ErrInvalidCredentials
		}
		return "", dto.UserResponse{}, fmt.Errorf("database error: %w", err)
	}

	if err := ComparePassword(hashedPassword, req.Password); err != nil {
		slog.Warn("Failed login attempt: password mismatch", "username", req.Username)
		return "", dto.UserResponse{}, ErrInvalidCredentials
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", dto.UserResponse{}, fmt.Errorf("เกิดข้อผิดพลาดในการสร้างโทเค็น: คีย์ลับไม่ถูกตั้งค่า")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", dto.UserResponse{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return t, user, nil
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}