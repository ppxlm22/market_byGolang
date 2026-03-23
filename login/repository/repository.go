package repository

import (
	"go_shopmarket/database"
	"go_shopmarket/login/dto"
)

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetUserByUsername(username string) (dto.UserResponse, string, error) {
	var user dto.UserResponse
	var passwordHash string

	query := `SELECT id, username, password_hash, role FROM users WHERE username = $1`
	err := database.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &passwordHash, &user.Role)
	if err != nil {
		return dto.UserResponse{}, "", err
	}
	return user, passwordHash, nil
}
