package main

import (
	"github.com/gin-gonic/gin"
	"github.com/theycallmesabb/echo/internal/db"
)

func main() {
	r := gin.Default()

	// Initialize the database
	db.InitDB()

	// Example route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run()
}
