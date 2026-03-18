package main

import (
	"go_shopmarket/config"
	"go_shopmarket/register"
	"log"
	"go_shopmarket/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	_ = config.LoadConfig()
	database.ConnectDB()

	userRepo := register.NewRepository()      
	userService := register.NewService(userRepo)
	userHandler := register.NewHandler(userService)

	app := fiber.New()

	app.Post("/register", userHandler.Register_Service)
	log.Fatal(app.Listen(":5000"))

}
