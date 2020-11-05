package controllers

import (
	"kwanjai/helpers"
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type passwordUpdate struct {
	OldPassword  string `json:"old_password" binding:"required,min=8"`
	NewPassword1 string `json:"new_password1" binding:"required,min=8"`
	NewPassword2 string `json:"new_password2" binding:"required,min=8"`
}

// Login endpoint
func Login() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		login := new(models.LoginCredential)
		err := ginContext.ShouldBindJSON(login)
		var status int
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		status, username := login.Login()
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
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		registerInfo.HashPassword()
		status, message, user = registerInfo.Register()
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
		extractedToken := strings.Split(ginContext.Request.Header.Get("Authorization"), "Bearer ")
		if len(extractedToken) != 2 {
			token.AccessToken = ""
		} else if extractedToken[0] != "Bearer " {
			token.AccessToken = ""
		} else {
			token.AccessToken = extractedToken[1]
		}
		if token.RefreshToken == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "No refresh token provied."})
			return
		}
		accessPassed := make(chan bool)
		accessTokenID := make(chan string)
		refreshPassed := make(chan bool)
		refreshTokenID := make(chan string)
		go logout.Verify(token.AccessToken, "access", accessPassed, accessTokenID)
		go logout.Verify(token.RefreshToken, "refresh", refreshPassed, refreshTokenID)
		passed := true == <-accessPassed && true == <-refreshPassed
		if !passed {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "Token verification failed."})
			return
		}
		go libraries.FirestoreDelete("tokens", <-accessTokenID)
		go libraries.FirestoreDelete("tokens", <-refreshTokenID)

		ginContext.JSON(200, gin.H{"message": "User logged out successfully."})
	}
}

// RefreshToken endpiont
func RefreshToken() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		token := new(libraries.Token)
		ginContext.ShouldBind(token)
		extractedToken := strings.Split(ginContext.Request.Header.Get("Authorization"), "Bearer ")
		if len(extractedToken) != 2 {
			token.AccessToken = ""
		} else if extractedToken[0] != "Bearer " {
			token.AccessToken = ""
		} else {
			token.AccessToken = extractedToken[1]
		}
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
		_, user, tokenID, err := libraries.VerifyToken(token.AccessToken, "access") // if token is expried here, it's got delete.
		if user != "anonymous" && err == nil {                                      // user != "anonymous" means token is still valid.
			_, err = libraries.FirestoreDelete("tokens", tokenID)
			if err != nil {
				ginContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		}
		newToken, err := libraries.CreateToken("access", user)
		token.AccessToken = newToken
		ginContext.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

// TokenVerification endpiont
func TokenVerification() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		if username != "anonymous" {
			ginContext.Status(http.StatusOK)
		} else {
			ginContext.Status(http.StatusUnauthorized)
		}
	}
}

// PasswordUpdate endpoint
func PasswordUpdate() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		passwordForm := new(passwordUpdate)
		if err := ginContext.ShouldBindJSON(passwordForm); err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		if passwordForm.NewPassword1 != passwordForm.NewPassword2 {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Password confrimation failed."})
		}
		username := helpers.GetUsername(ginContext)
		newPassword, _ := libraries.HashPassword(passwordForm.NewPassword1)
		libraries.FirestoreUpdateField("users", username, "HashedPassword", newPassword)
		ginContext.JSON(http.StatusOK, gin.H{"message": "Password updated."})
	}
}
