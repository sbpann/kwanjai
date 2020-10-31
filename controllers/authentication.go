package controllers

import (
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Login endpoint
func Login() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		login := new(models.LoginCredential)
		err := ginContext.ShouldBindJSON(login)
		var status int
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Login form is not valid."})
			return
		}
		status, username := models.Login(login)
		if status != http.StatusOK {
			ginContext.JSON(status, gin.H{"message": username})
			// If status != 200, error message is returned instead of username.
			return
		}
		token := new(libraries.Token)
		token.Initialize(username)
		ginContext.JSON(status, gin.H{"user": username, "token": token})
	}
}

// Register endpoint
func Register() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		registerInfo := new(models.User)
		// Keep in mind.
		// if content type is not provided ShouldBind is ShouldBindForm.
		err := ginContext.ShouldBindJSON(registerInfo)
		var status int
		var message string
		var user *models.User
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Registration form is not valid."})
			return
		}
		registerInfo.HashPassword()
		status, message, user = models.Register(registerInfo)
		if status != http.StatusOK {
			ginContext.JSON(status, gin.H{"message": message})
			return
		}
		token := new(libraries.Token)
		token.Initialize(user.Username)
		if registerInfo.Username == "test" || registerInfo.Email == "test@example.com" {
			ginContext.JSON(status, gin.H{
				"message": message,
				"user":    user,
				"token":   token,
				"warning": "You have just registered with the username (test) or the email (test@example.com) which is going to be delete eventually. Please avoid using those names.",
			})
			return
		}
		ginContext.JSON(status, gin.H{
			"message": message,
			"user":    user,
			"token":   token,
		})

	}
}

// Logout endpoint
func Logout() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		user, exist := ginContext.Get("user")
		if !exist {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "No user found in backend context."})
		}
		var (
			passed           bool
			accessPassed     bool
			refreshPassed    bool
			accessTokenUUID  string
			refreshTokenUUID string
		)
		logout := new(models.LogoutData)
		token := new(libraries.Token)
		ginContext.ShouldBind(token)
		token.AccessToken = ginContext.Request.Header.Get("Authorization")
		go logout.Verify(token.AccessToken, "access")
		go logout.Verify(token.RefreshToken, "refresh")
		timeout := time.Now().Add(time.Second * 4)
		timer := time.Now()
		for !passed && timer.Before(timeout) {
			accessPassed = logout.AccessPassed
			refreshPassed = logout.RefreshPassed
			accessTokenUUID = logout.AccessTokenUUID
			refreshTokenUUID = logout.RefreshTokenUUID
			passed = accessPassed == true && refreshPassed == true
			timer = time.Now()
		}
		if !passed {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "token verification failed."})
			return
		}
		go libraries.DeleteToken(user.(string), accessTokenUUID)
		go libraries.DeleteToken(user.(string), refreshTokenUUID)

		ginContext.JSON(200, gin.H{"message": "User logged out successfully."})
	}
}

// RefreshToken endpiont
func RefreshToken() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		token := new(libraries.Token)
		ginContext.ShouldBind(token)
		token.AccessToken = ginContext.Request.Header.Get("Authorization")
		if token.RefreshToken == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "No refresh token provied."})
			return
		}
		passed, user, _, err := libraries.VerifyToken(token.RefreshToken, "refresh")
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if !passed && user != "anonymous" {
			var status int
			if err.Error() == "user not found" {
				status = http.StatusNotFound
			} else {
				status = http.StatusInternalServerError
			}
			ginContext.JSON(status, gin.H{"message": err.Error()})
			return
		}
		_, user, tokenUUID, err := libraries.VerifyToken(token.AccessToken, "access") // if token is expried here. it got delete.
		if user != "anonymous" && err == nil {                                        // which means token is still valid.
			err = libraries.DeleteToken(user, tokenUUID)
			if err != nil {
				ginContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		}
		newToken, err := libraries.CreateToken("access", user)
		token.AccessToken = newToken
		ginContext.JSON(http.StatusOK, gin.H{
			"user":  user,
			"token": token,
		})
	}
}
