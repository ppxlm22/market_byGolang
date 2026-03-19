package register

import "go_shopmarket/database"

type Repository interface {
	CheckUserExists(username, email string) (bool, error)
	Register(req registerDB) error

}
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

func (r *repository) Register(req registerDB) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`

	_, err := database.DB.Exec(query, req.Username, req.Email, req.Password)
	if err != nil {
		return err
	}
	return nil
}
