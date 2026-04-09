package repository

import (
	"database/sql"
	"go_shopmarket/database"
	"go_shopmarket/register/dto"
)

type repository struct{
	DB *sql.DB
}

func NewRepository(DB *sql.DB) Repository {
	return &repository{
		DB: DB,
	}
}
func (r *repository) CheckUserExists(username string, email string) (bool, error) {
	var exists bool
	
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)`
	
	err := r.DB.QueryRow(query, username, email).Scan(&exists)
	if err != nil {
		return false, err 
	}
	return exists, nil 
}

func (r *repository) Register(req dto.RegisterRequest) (*dto.RegisterDB, error) {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`

	_, err := database.DB.Exec(query, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &dto.RegisterDB{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}, nil
}
