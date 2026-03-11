package shopHandler

import (
	"fmt"
	"go_shopmarket/database"
	"log"
)
func AddCategory (name string) {
	query := `INSERT INTO categories (name) VALUE ($1) RETURNING id`

	var lastInsertID int 
	err := database.DB.QueryRow(query, name)

	if err != nil {
		log.Println("Insert category Error", err)
		return
	}
	fmt.Printf("add category '%s' success (ID: %d )", name, lastInsertID)
}