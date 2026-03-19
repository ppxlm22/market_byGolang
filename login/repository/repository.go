package repository

import (
	"go_shopmarket/database"
)
type repository struct{}
func NewRepository() Repository{
	return &repository{}
}

func (r *repository) GetUserByUsername(username string) (int, string, error) {
	var id int
	var passwordHash string

	query := `SELECT id, password_hash FROM users WHERE username = $1`
	err := database.DB.QueryRow(query, username).Scan(&id, &passwordHash)
	if err != nil {
		return 0, "", err
	}
	return id, passwordHash, nil

}
