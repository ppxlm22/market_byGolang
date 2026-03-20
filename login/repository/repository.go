package repository

import (
	"go_shopmarket/database"
)
type repository struct{}
func NewRepository() Repository{
	return &repository{}
}

func (r *repository) GetUserByUsername(username string) (int, string, string, error) {
	var id int
	var passwordHash string
	var role string

	query := `SELECT id, password_hash, role FROM users WHERE username = $1`
	err := database.DB.QueryRow(query, username).Scan(&id, &passwordHash, &role)
	if err != nil {
		return 0, "", "", err
	}
	return id, passwordHash, role, nil
}
