package main

import (
	"github.com/gin-gonic/gin"
	"github.com/theycallmesabb/echo/internal/auth"
	"github.com/theycallmesabb/echo/internal/db"
)

func main() {
	r := gin.Default()

	// Initialize the database
	db.InitDB()
	auth.InitGoogleauth()
	r.GET("/login", auth.HandleGoogleLogin)
	r.GET("/callback", auth.HandleGoogleCallback)

	r.Run()
}
