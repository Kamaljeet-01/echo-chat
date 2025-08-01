package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGoogleLogin(c *gin.Context) {
	url := GetGoogleOAuthConfig().AuthCodeURL("random-state-token")
	c.Redirect(http.StatusTemporaryRedirect, url)
}
