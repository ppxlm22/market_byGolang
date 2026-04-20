package repository

import (
	"database/sql"
	"go_shopmarket/register/dto"
	"errors"
)

var ErrDBQuery = errors.New("database query error")
var ErrUserAlreadyExists = errors.New("user already exists in database")

type repository struct{
	DB *sql.DB
}

func NewRepository(DB *sql.DB) Repository {
	return &repository{
		DB: DB,
	}
}
//
func (r *repository) CheckUserExists(username string, email string) (bool, error) {
	var exists bool
	
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)`
	
	err := r.DB.QueryRow(query, username, email).Scan(&exists)
	if err != nil {
		return false, ErrDBQuery 
	}
	return exists, nil 
}

func (r *repository) Register(req dto.RegisterRequest) (*dto.RegisterDB, error) {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`

	_, err := r.DB.Exec(query, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, ErrDBQuery
	}
	return &dto.RegisterDB{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}, nil
}
