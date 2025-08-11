package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		//   First check Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			if !strings.HasPrefix(authHeader, "Bearer ") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Bearer prefix is missing",
				})
				return
			}
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		//   If no header, check cookie
		if token == "" {
			cookieToken, err := c.Cookie("idtoken")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "No token found in header or cookie",
				})
				return
			}
			token = cookieToken
		}

		//   Validate ID token (audience "" means skip audience check)
		payload, err := idtoken.Validate(context.Background(), token, "")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		//  Store user email in context also this is a type assertion
		email, ok := payload.Claims["email"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Email not found in token",
			})
			return
		}

		c.Set("user", email)
		c.Next()
	}
}
