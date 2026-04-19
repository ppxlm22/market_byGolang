package repository

import (
	"database/sql"
	"go_shopmarket/login/dto"
	"errors"
	"fmt"
)
var ErrDBQuery = errors.New("database query error")
var ErrUserNotFound = errors.New("user not found")

type repository struct{
	DB *sql.DB
}

func NewRepository(DB *sql.DB) Repository {
	return &repository{
		DB: DB,
	}
}

func (r *repository) GetUserByUsername(username string) (dto.UserResponse, string, error) {
	var user dto.UserResponse
	var passwordHash string

	query := `SELECT id, username, password_hash, role FROM users WHERE username = $1`
	
	err := r.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &passwordHash, &user.Role)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.UserResponse{}, "", ErrUserNotFound
		}
		return dto.UserResponse{}, "", fmt.Errorf("%w: %v", ErrDBQuery, err)
	}
	
	return user, passwordHash, nil
}