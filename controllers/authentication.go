package controllers

import (
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Login function returns status of password validiation as booleans
func Login(ginContext *gin.Context) {
	login := new(models.LoginCredential)
	err := ginContext.ShouldBindJSON(login)
	var status int
	var message string
	if err != nil {
		status, message = http.StatusBadRequest, "Login form is not valid."
		ginContext.JSON(status, gin.H{"message": message})
		return
	}
	status, username := models.Login(login)
	if status != http.StatusOK {
		ginContext.JSON(status, gin.H{"message": username})
		return
	}
	token := new(libraries.Token)
	token.Initialize(username)
	ginContext.JSON(status, gin.H{"user": username, "token": token})
}

// Register new user
func Register(ginContext *gin.Context) {
	registerInfo := new(models.User)
	err := ginContext.ShouldBind(registerInfo)
	var status int
	var message string
	var user *models.User
	if err != nil {
		status, message = http.StatusBadRequest, "Registration form is not valid."
		ginContext.JSON(status, gin.H{"message": message})
		return
	}
	registerInfo.HashPassword()
	status, message, user = models.Register(registerInfo)
	if status != http.StatusCreated {
		ginContext.JSON(status, gin.H{"message": message})
		return
	}
	token := new(libraries.Token)
	token.Initialize(user.Username)
	ginContext.JSON(status, gin.H{
		"message": message,
		"user":    user,
		"token":   token,
	})

}

func Logout(ginContext *gin.Context) {
	var (
		pass             bool
		accessPass       bool
		refreshPass      bool
		user             string
		accessTokenUUID  string
		refreshTokenUUID string
		err              error
	)
	logout := new(models.LogoutData)
	token := new(libraries.Token)
	err = ginContext.ShouldBind(&token)
	token.AccessToken = ginContext.Request.Header.Get("Authorization")
	go logout.Verify(token.AccessToken, "access")
	go logout.Verify(token.RefreshToken, "refresh")
	timeout := time.Now().Add(time.Second * 10)
	timer := time.Now()
	for !pass && !timer.Equal(timeout) {
		accessPass = logout.AccessPass
		refreshPass = logout.RefreshPass
		user = logout.User
		accessTokenUUID = logout.AccessTokenUUID
		refreshTokenUUID = logout.RefreshTokenUUID
		err = logout.Err
		pass = accessPass == true && refreshPass == true
		timer = time.Now()
	}
	if !pass {
		ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "token verification failed."})
		return
	}
	if !pass && user != "anonymous" {
		var status int
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		ginContext.JSON(status, gin.H{"message": err.Error()})
		return
	} else if !pass && user == "anonymous" {
		ginContext.JSON(200, gin.H{"message": err.Error()})
		return
	}

	if pass {
		go libraries.DeleteToken(user, accessTokenUUID)
		go libraries.DeleteToken(user, refreshTokenUUID)
	}

	ginContext.JSON(200, gin.H{"message": "logout successfully"})
}

func AuthenticateTest(ginContext *gin.Context) {
	pass, user, _, err := libraries.VerifyToken(ginContext.Request.Header.Get("Authorization"), "access")
	if !pass && user != "anonymous" {
		var status int
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		ginContext.JSON(status, gin.H{"message": err.Error()})
		return
	}
	ginContext.JSON(200, gin.H{"pass": pass, "user": user})
}
