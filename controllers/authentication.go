package controllers

import (
	"gin-sandbox/libraries"
	"gin-sandbox/models"

	"github.com/gin-gonic/gin"
)

// Login function returns status of password validiation as booleans
func Login(c *gin.Context) {
	user := new(models.User)
	c.BindJSON(user)
	status, detail := user.Login()
	if status != 200 {
		c.JSON(status, gin.H{
			"detail": detail,
		})
	} else {
		c.JSON(status, gin.H{
			"detail": detail,
		})
	}
}

// Register new user
func Register(c *gin.Context) {
	user := new(models.User)
	c.BindJSON(&user)
	hashedpassword, _ := libraries.HashPassword(user.Password)
	user.Password = hashedpassword
	status, detail := user.CreateUser()
	c.JSON(status, gin.H{
		"detail": detail,
	})
}
