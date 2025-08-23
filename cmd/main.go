package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"echo/internal/auth"
	"echo/internal/chat"
	"echo/internal/db"
	"echo/internal/user"
	"echo/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// POST /messages
func SendMessage(c *gin.Context) {
	var msg user.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	senderEmail := c.MustGet("email").(string)
	msg.SenderEmail = senderEmail
	msg.CreatedAt = time.Now()

	if err := db.CreateMessage(&msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save message"})
		return
	}

	msgBytes, err := json.Marshal(msg)
	if err == nil {
		go chat.ForwardMessage(msg.ReceiverEmail, msgBytes)
	}

	c.JSON(http.StatusOK, gin.H{"status": "message sent"})
}

// GET /messages/:receiver
func GetMessages(c *gin.Context) {
	sender := c.MustGet("email").(string)
	receiver := c.Param("receiver")

	msgs, err := db.GetMessage(sender, receiver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, msgs)
}

func main() {
	fmt.Println(" Starting server...")

	// Init
	fmt.Println(" Initializing DB...")
	db.InitDB()
	fmt.Println(" DB Init done.")

	fmt.Println(" Initializing Google Auth...")
	auth.InitGoogleauth()
	fmt.Println(" Google Auth Init done.")

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Load HTML templates
	r.LoadHTMLGlob("*.html")

	// Serve the main HTML page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Public Routes
	fmt.Println(" Registering public routes...")
	r.GET("/login", func(c *gin.Context) {
		fmt.Println("➡ /login handler triggered")
		auth.HandleGoogleLogin(c)
	})
	r.GET("/callback", func(c *gin.Context) {
		fmt.Println("➡ /callback handler triggered")
		auth.HandleGoogleCallback(c)
	})

	// Use your cookie verification middleware for all routes in this group
	authRoutes := r.Group("/", middleware.Verifycookie())

	authRoutes.GET("/me", func(c *gin.Context) {
		email := c.MustGet("email").(string)
		c.JSON(200, gin.H{"message": "Hello!", "email": email})
	})

	// Message Routes
	authRoutes.POST("/messages", SendMessage)
	authRoutes.GET("/messages/:receiver", GetMessages)

	//  WebSocket route
	authRoutes.GET("/ws", func(c *gin.Context) {
		email, exists := c.Get("email")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		chat.HandleWebSocket(c.Writer, c.Request, email.(string))
	})

	fmt.Println(" Server is running on :8080")
	r.Run()
}
