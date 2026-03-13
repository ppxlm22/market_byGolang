package main

import (
	"github.com/gofiber/fiber/v2"
    "github.com/joho/godotenv"
    "go_shopmarket/database"
    "go_shopmarket/handlers"
    "log"
	"fmt"
)

func main() {
	if err := godotenv.Load();
	err != nil {
		log.Println("Warning: .env file not found")
	}
	fmt.Println("CONNECT ENV")
    database.ConnectDB()
    app := fiber.New()

	app.Get("/product/:id", shopHandler.Getproduct)
	app.Get("/products", shopHandler.Getallproducts)
	log.Fatal(app.Listen(":5000"))

}