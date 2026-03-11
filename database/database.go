package database 

import (
	"fmt"
	"log"
	"os"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)
var DB *sql.DB

func ConnectDB(){
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}	
	DB = db
	fmt.Println("CONNECT SUCCESSFULLY") 
}