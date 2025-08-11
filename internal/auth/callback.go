package auth

import (
	"context"
	"fmt"
	"net/http"

	"echo/internal/db"
	"echo/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"google.golang.org/api/oauth2/v1"
)

func HandleGoogleCallback(c *gin.Context) {
	code := c.Query("code")
	fmt.Println("Received code:", code)

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

	User := &user.Chatuser{
		Name:  userinfo.Name,
		Email: userinfo.Email,
	}
	check, err := db.Checkuser(User.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong" + err.Error(),
		})
		return
	}
	if check {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "User already exists",
		})
		return
	}

	err = db.Create(User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
		})
		return
	}
	idtoken, ok := token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ID Token not found",
		})
		return
	}
	c.SetCookie(
		"idtoken", // Cookie name
		idtoken,   // Value (the token)
		360000,    // MaxAge in seconds (100 hour)
		"/",       // Path
		"",        // Domain (empty = current domain
		true,      // Secure (only sent over HTTPS)
		true,      // HttpOnly (not accessible via JS)

	)

	c.JSON(http.StatusOK, gin.H{
		"email": userinfo.Email,
		"name":  userinfo.Name,

		// "id_token": idtoken, // remove if you don’t want to expose it to JS
	})
}