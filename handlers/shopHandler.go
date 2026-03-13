package shopHandler

import (
	"database/sql"
	"go_shopmarket/database"
	_ "go_shopmarket/database"
	"go_shopmarket/models"

	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Getproduct(c *fiber.Ctx) error {
	id := c.Params("id")
	fmt.Println(id)

	var p models.Products

	query := `SELECT *
	FROM public.products WHERE id = $1;`

	err := database.DB.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.Create_at)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString("Product not found")
		}
		return err
	}

	return c.JSON(&p)
}
func Getallproducts(c *fiber.Ctx) error {
	var p models.Products
	var proDucts []models.Products
	
	query := `SELECT * FROM public.products;`
	rows, _ := database.DB.Query(query)

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.Create_at); err != nil {
			return err
		}
		fmt.Println(p)
		proDucts = append(proDucts,p)
	}

	return c.JSON(proDucts)

}
func Createproduct(c *fiber.Ctx) error {
	var p models.Products
	query := `INSERT 
				INTO public.products(
				id, name, price, stock, category_id, create_at)
				VALUES (?, ?, ?, ?, ?, ?);`

	
}