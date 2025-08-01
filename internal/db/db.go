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

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DATABASE: %v", err)

	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Connected to DB & Migrated successfully")
}
