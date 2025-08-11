package main

import (
	"fmt"

	"echo/internal/auth"
	"echo/internal/db"
	"echo/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("🚀 Starting server...")

	// Init
	fmt.Println("📦 Initializing DB...")
	db.InitDB()
	fmt.Println("✅ DB Init done.")

	fmt.Println("🔑 Initializing Google Auth...")
	auth.InitGoogleauth()
	fmt.Println("✅ Google Auth Init done.")

	r := gin.Default()

	// Public Routes
	fmt.Println("📍 Registering public routes...")
	r.GET("/login", func(c *gin.Context) {
		fmt.Println("➡ /login handler triggered")
		auth.HandleGoogleLogin(c)
	})
	r.GET("/callback", func(c *gin.Context) {
		fmt.Println("➡ /callback handler triggered")
		auth.HandleGoogleCallback(c)
	})

	// Protected Routes
	fmt.Println("📍 Registering protected routes...")
	authRoutes := r.Group("/", middleware.AuthMiddleware())
	authRoutes.GET("/me", func(c *gin.Context) {
		fmt.Println("➡ /me handler triggered")
		email := c.MustGet("user").(string)
		c.JSON(200, gin.H{"message": "Hello!", "email": email})
	})

	r.GET("/test", func(ctx *gin.Context) {
		fmt.Println("➡ /test handler triggered")
		c, err := ctx.Cookie("idtoken")
		if err != nil {
			fmt.Println("⚠ No cookie found")
			return
		}
		fmt.Println("🍪 Cookie from browser:", c)
	})

	fmt.Println("✅ Server is running on :8080")
	r.Run()
}
