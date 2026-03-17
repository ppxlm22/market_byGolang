package register

import "go_shopmarket/database"

func Register(req registerDB) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`

	_, err := database.DB.Exec(query, req.Username, req.Email, req.Password)
	if err != nil {
		return err
	}
	return nil
}