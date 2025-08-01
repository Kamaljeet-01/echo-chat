package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/oauth2/v1"
)

func HandleGoogleCallback(c *gin.Context) {
	code := c.Query("code")

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	service, err := oauth2.New(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create oauth2 service"})
		return
	}

	userinfo, err := service.Userinfo.Get().Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":     userinfo.Email,
		"name":      userinfo.Name,
		"verified":  userinfo.VerifiedEmail,
		"google_id": userinfo.Id,
	})
}
