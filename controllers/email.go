package controllers

import (
	"kwanjai/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// VerifyEmail endpoint
func VerifyEmail() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		verificationEmail := new(models.VerificationEmail)
		ginContext.ShouldBind(verificationEmail)
		verificationEmail.UUID = ginContext.Param("UUID")
		status, message := verificationEmail.Verify()
		ginContext.JSON(status, gin.H{"message": message})
	}
}

// ResendVerifyEmail endpoint
func ResendVerifyEmail() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		verificationEmail := new(models.VerificationEmail)
		ginContext.ShouldBind(verificationEmail)
		user := new(models.User)
		user.Email = verificationEmail.Email
		status, message, user := user.Finduser()
		if user == nil {
			ginContext.JSON(status, gin.H{"message": message})
			return
		}
		if user.IsVerified {
			ginContext.JSON(http.StatusOK, gin.H{"message": "The user is already verified."})
			return
		}
		status, message = user.SendVerificationEmail()
		ginContext.JSON(status, gin.H{"message": message})
	}
}
