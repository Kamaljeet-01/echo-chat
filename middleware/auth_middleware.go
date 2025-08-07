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

		//get the authorization header from request

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			//if header is missing black the request with 401 unauthorized
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Authorization header is missing",
			})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		//if beared prefix is not there then its not a valid format

		if token == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error ": "Bearer prefix is missing",
			})
		}

		payload, err := idtoken.Validate(context.Background(), token, "")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}
		// if valid, set the user email in context so it can be used later

		c.Set("user", payload.Claims["email"])
		c.Next()
	}
}
