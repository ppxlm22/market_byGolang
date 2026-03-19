package main

import (
	"go_shopmarket/config"
	"go_shopmarket/register"
	"log"
	"go_shopmarket/database"
	loginRepo "go_shopmarket/login/repository"
	loginSvc  "go_shopmarket/login/service"
	loginHdl  "go_shopmarket/login/handler"
	

	"github.com/gofiber/fiber/v2"
)

func main() {
	_ = config.LoadConfig()
	database.ConnectDB()

	userRepo := register.NewRepository()      
	userService := register.NewService(userRepo)
	userHandler := register.NewHandler(userService)

	loginRepo := loginRepo.NewRepository()
	loginService := loginSvc.NewService(loginRepo)
	loginHandler := loginHdl.NewHandler(loginService)

	app := fiber.New()

	app.Post("/login", loginHandler.Login_Service)
	app.Post("/register", userHandler.Register_Service)
	log.Fatal(app.Listen(":5000"))

}
