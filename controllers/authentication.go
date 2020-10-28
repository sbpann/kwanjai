package controllers

import (
	"gin-sandbox/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login function returns status of password validiation as booleans
func Login(c *gin.Context) {
	login := new(models.LoginCredential)
	err := c.ShouldBindJSON(login)
	var status int
	var detail string
	if err != nil {
		status, detail = http.StatusBadRequest, "Login form is not valid."
	} else {
		status, detail = models.Login(login)
	}
	c.JSON(status, gin.H{
		"detail": detail,
	})
}

// Register new user
func Register(c *gin.Context) {
	user := new(models.User)
	err := c.ShouldBindJSON(user)
	var status int
	var detail string
	if err != nil {
		status, detail = http.StatusBadRequest, "Registration form is not valid."
	} else {
		user.HashPassword()
		status, detail = models.Register(user)
	}
	c.JSON(status, gin.H{
		"detail": detail,
	})
}
