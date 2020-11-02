package controllers

import (
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"

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
		ginContext.JSON(status, gin.H{"token": token})
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
		if registerInfo.Username == "test1" ||
			registerInfo.Email == "test1@example.com" ||
			registerInfo.Username == "test2" ||
			registerInfo.Email == "test2@example.com" {
			ginContext.JSON(status, gin.H{
				"message": message,
				"token":   token,
				"warning": "You have just registered with the username (test) or the email (test@example.com) which is going to be delete eventually. Please avoid using those names.",
			})
			return
		}
		ginContext.JSON(status, gin.H{
			"message": message,
			"token":   token,
		})

	}
}

// Logout endpoint
func Logout() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		logout := new(models.LogoutData)
		token := new(libraries.Token)
		ginContext.ShouldBindJSON(token)
		// Todo: add Bearer prefix
		token.AccessToken = ginContext.Request.Header.Get("Authorization")
		if token.RefreshToken == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "No refresh token provied."})
			return
		}
		accessPassed := make(chan bool)
		accessTokenUUID := make(chan string)
		refreshPassed := make(chan bool)
		refreshTokenUUID := make(chan string)
		go logout.Verify(token.AccessToken, "access", accessPassed, accessTokenUUID)
		go logout.Verify(token.RefreshToken, "refresh", refreshPassed, refreshTokenUUID)
		passed := true == <-accessPassed && true == <-refreshPassed
		if !passed {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "Token verification failed."})
			return
		}
		go libraries.FirestoreDelete("tokenUUID", <-accessTokenUUID)
		go libraries.FirestoreDelete("tokenUUID", <-refreshTokenUUID)

		ginContext.JSON(200, gin.H{"message": "User logged out successfully."})
	}
}

// RefreshToken endpiont
func RefreshToken() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		token := new(libraries.Token)
		ginContext.ShouldBind(token)
		// Todo: add Bearer prefix
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
		_, user, tokenUUID, err := libraries.VerifyToken(token.AccessToken, "access") // if token is expried here, it's got delete.
		if user != "anonymous" && err == nil {                                        // user != "anonymous" means token is still valid.
			_, err = libraries.FirestoreDelete("tokenUUID", tokenUUID)
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
