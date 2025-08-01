package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	Id      uint
	Email   string
	Name    string
	Picture string
}

var db *gorm.DB

func InitDB() {
	errs := godotenv.Load()
	if errs != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_URL")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DATABASE: %v", err)

	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Connected to DB & Migrated successfully")
}
