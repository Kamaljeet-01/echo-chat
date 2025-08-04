package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/theycallmesabb/echo/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	errs := godotenv.Load()
	if errs != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_URL")
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DATABASE: %v", err)

	}
	err = Db.AutoMigrate(&user.Chatuser{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Connected to DB & Migrated successfully")

}

func Create(ch any) error {

	err := Db.Create(ch).Error
	if err != nil {
		return err
	}
	return nil
}

func Find(ch any) (*user.Chatuser, error) {
	var user user.Chatuser
	err := Db.Raw("SELECT * FROM Chat_user").Scan(&user).Error
	if err != nil {
		return &user, err
	}
	return &user, nil

}

func Checkuser(email string) (bool, error) {
	var user user.Chatuser
	err := Db.Raw("SELECT * FROM chat_user where email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, err
		}
		return false, fmt.Errorf("err in checking for user, err: %v", err)
	}
	fmt.Println("err:", err)
	return true, nil

}
