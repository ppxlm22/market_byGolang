package main

import (
	"go_shopmarket/database" 
	"go_shopmarket/handler"
)

func main() {
	database.ConnectDB()

	shopHandler.AddCategory("Electronics")
}