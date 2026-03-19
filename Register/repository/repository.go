package repository

import (
	"go_shopmarket/database"
	"go_shopmarket/register/dto"
)

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}
func (r *repository) CheckUserExists(username string, email string) (bool, error) {
	var exists bool
	
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)`
	
	err := database.DB.QueryRow(query, username, email).Scan(&exists)
	if err != nil {
		return false, err 
	}
	return exists, nil 
}

func (r *repository) Register(req dto.RegisterRequest) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`

	_, err := database.DB.Exec(query, req.Username, req.Email, req.Password)
	if err != nil {
		return err
	}
	return nil
}
