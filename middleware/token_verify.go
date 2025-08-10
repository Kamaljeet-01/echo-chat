package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

func Verifycookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("idtoken")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No cookie found"})
			c.Abort()
			return
		}

		cid := os.Getenv("GOOGLE_CLIENT_ID")
		payload, err := idtoken.Validate(context.Background(), cookie, cid)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store email and name in context for later use
		c.Set("email", payload.Claims["email"])
		c.Set("name", payload.Claims["name"])

		c.Next()
	}
}
