package shopHandler

import (
	"database/sql"
	"go_shopmarket/database"
	_ "go_shopmarket/database"
	"go_shopmarket/models"
	"strconv"

	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)
func Register(c *fiber.Ctx) error {
	type RegisterRequest struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
	var requ
}	

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
func GetAllproducts(c *fiber.Ctx) error {
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
		proDucts = append(proDucts, p)
	}

	return c.JSON(proDucts)

}
func CreateProduct(c *fiber.Ctx) error {
	fmt.Println("CreateProduct")
	p := new(models.Products)
	if err := c.BodyParser(p); err != nil {
		return err
	}
	_, err := database.DB.Exec("INSERT INTO public.products (name, price, stock, category_id ) VALUES ($1, $2, $3, $4)", p.Name, p.Price, p.Stock, p.CategoryID)
	if err != nil {
		return err
	}
	return c.JSON(p)

}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	p := new(models.Products)
	if err := c.BodyParser(p); err != nil {
		return err
	}
	_, err := database.DB.Exec("UPDATE public.products SET name = $1, price = $2,stock = $3 ,category_id = $4 WHERE id = $5", p.Name, p.Price, p.Stock, p.CategoryID, id)
	if err != nil {
		return err
	}
	p.ID, _ = strconv.Atoi(id)
	return c.JSON(p)
}
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := database.DB.Exec("DELETE FROM public.products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
