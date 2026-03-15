package main

import (
	"fmt"
	"go_shopmarket/database"
	shopHandler "go_shopmarket/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}
	fmt.Println("CONNECT ENV")
	database.ConnectDB()
	app := fiber.New()

	app.Post("/register", shopHandler.Register)
	app.Post("/login", shopHandler.Login)

	app.Post("/products", shopHandler.CreateProduct)
	app.Get("/product/:id", shopHandler.Getproduct)
	app.Get("/products", shopHandler.GetAllproducts)

	app.Put("/product/:id", shopHandler.UpdateProduct)
	app.Delete("/product/:id", shopHandler.DeleteProduct)
	log.Fatal(app.Listen(":5000"))

}
