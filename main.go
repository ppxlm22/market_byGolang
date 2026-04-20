package main

import (
	"errors"
	"go_shopmarket/config"
	"log/slog"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go_shopmarket/database"
	"log"
	"go_shopmarket/middleware"
	"github.com/gofiber/fiber/v2"

	registerHdl "go_shopmarket/register/handler"
	register "go_shopmarket/register/repository"
	registerSvc "go_shopmarket/register/service"

	loginHdl "go_shopmarket/login/handler"
	loginRepo "go_shopmarket/login/repository"
	loginSvc "go_shopmarket/login/service"

	productHdl "go_shopmarket/products/handler"
	productRepo "go_shopmarket/products/repository"
	productSvc "go_shopmarket/products/service"
	
)

func main() {
	
	_ = config.LoadConfig()

	DB := database.ConnectDB()
	// MysqlDB := database.ConnectMysql()

	userRepo := register.NewRepository(DB)
	userService := registerSvc.NewService(userRepo)
	userHandler := registerHdl.NewHandler(userService)

	loginRepo := loginRepo.NewRepository(DB)
	loginService := loginSvc.NewService(loginRepo)
	loginHandler := loginHdl.NewHandler(loginService)

	productRepo := productRepo.NewRepository(DB)
	productService := productSvc.NewService(productRepo)
	productHandler := productHdl.NewHandler(productService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func (c *fiber.Ctx, err error) error  {
			code := fiber.StatusInternalServerError
			message := "เกิดข้อผิดพลาดที่เซิฟเวอร์"

			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr){
				code = fiberErr.Code
				message = fiberErr.Message
			}else{
				slog.Error("Unhanler Server Error",
					"error", err,
					"path", c.Path(),
					"method", c.Method(),
				)
			}
			return c.Status(code).JSON(fiber.Map{
				"error": message,
			})
		},
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	app.Post("/login", loginHandler.Login_Service)
	app.Post("/register", userHandler.Register_Service)
	app.Post("/products", middleware.Protected(),middleware.AdminOnly(), productHandler.CreateProduct)
	app.Get("/products", productHandler.GetAllProducts)
	app.Get("/products/:id", productHandler.GetProductByID)
	app.Put("/products/:id", middleware.Protected(),middleware.AdminOnly(), productHandler.UpdateProduct)
	app.Delete("/products/:id", middleware.Protected(),middleware.AdminOnly(), productHandler.DeleteProduct)
	app.Get("/categories/:id", productHandler.GetCategoryByID)
	app.Get("/categories", productHandler.GetAllCategories)
	app.Post("/checkout",middleware.Protected(), productHandler.Checkout_service)
	app.Static("/", "./public")
	log.Fatal(app.Listen(":5080"))

}
