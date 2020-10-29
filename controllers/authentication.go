package controllers

import (
	"fmt"
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login function returns status of password validiation as booleans
func Login(c *gin.Context) {
	login := new(models.LoginCredential)
	err := c.ShouldBindJSON(login)
	var status int
	var message string
	if err != nil {
		status, message = http.StatusBadRequest, "Login form is not valid."
		c.JSON(status, gin.H{"message": message})
		return
	}
	status, username := models.Login(login)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"message": username})
		return
	}
	token := new(libraries.Token)
	token.Initialize(username)
	c.JSON(status, gin.H{"user": username, "token": token})
}

// Register new user
func Register(c *gin.Context) {
	registerInfo := new(models.User)
	err := c.ShouldBind(registerInfo)
	var status int
	var message string
	var user *models.User
	fmt.Println(err)
	if err != nil {
		status, message = http.StatusBadRequest, "Registration form is not valid."
		c.JSON(status, gin.H{"message": message})
		return
	}
	registerInfo.HashPassword()
	status, message, user = models.Register(registerInfo)
	if status != http.StatusCreated {
		c.JSON(status, gin.H{"message": message})
		return
	}
	token := new(libraries.Token)
	token.Initialize(user.Username)
	c.JSON(status, gin.H{
		"message": message,
		"user":    user,
		"token":   token,
	})

}

func AuthenticateTest(c *gin.Context) {
	pass, user, err := libraries.VerifyToken(c.Request.Header.Get("Authorization"), "access")
	if !pass && user != "anonymous" {
		c.JSON(500, gin.H{"message": err.Error()})
	}
	c.JSON(200, gin.H{"pass": pass, "user": user, "error": err.Error()})
}
