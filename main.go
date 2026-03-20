package main

import (
	"go_shopmarket/config"
	
	"log"
	"go_shopmarket/database"
	register "go_shopmarket/register/repository"
	registerSvc "go_shopmarket/register/service"
	registerHdl "go_shopmarket/register/handler"

	loginRepo "go_shopmarket/login/repository"
	loginSvc  "go_shopmarket/login/service"
	loginHdl  "go_shopmarket/login/handler"

	productRepo "go_shopmarket/products/repository"
	productSvc  "go_shopmarket/products/service"
	productHdl  "go_shopmarket/products/handler"



	"github.com/gofiber/fiber/v2"
)

func main() {
	_ = config.LoadConfig()
	database.ConnectDB()

	userRepo := register.NewRepository()      
	userService := registerSvc.NewService(userRepo)
	userHandler := registerHdl.NewHandler(userService)

	loginRepo := loginRepo.NewRepository()
	loginService := loginSvc.NewService(loginRepo)
	loginHandler := loginHdl.NewHandler(loginService)

	productRepo := productRepo.NewRepository()
	productService := productSvc.NewService(productRepo)
	productHandler := productHdl.NewHandler(productService)

	app := fiber.New()

	app.Post("/login", loginHandler.Login_Service)
	app.Post("/register", userHandler.Register_Service)
	app.Post("/products", productHandler.CreateProduct)
	app.Get("/products", productHandler.GetAllProducts)
	app.Get("/product/:id", productHandler.GetProductByID)
	app.Put("/product/:id", productHandler.UpdateProduct)
	log.Fatal(app.Listen(":5000"))

}
