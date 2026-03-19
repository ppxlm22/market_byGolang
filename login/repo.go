package login

import (
	"go_shopmarket/database"
)

type Repository interface {
	GetUserByUsername(username string) (int, string, error)
}
type repository struct{}

func NewRepository() Repository{
	return &repository{}
}

func (r *repository) GetUserByUsername(username string) (int, string, error) {
	var id int
	var passwordHash string

	query := `SELECT password_hash FROM users WHERE username = $1`
	err := database.DB.QueryRow(query, username).Scan(&id, &passwordHash)
	if err != nil {
		return 0, "", err
	}
	return id, passwordHash, nil

}
