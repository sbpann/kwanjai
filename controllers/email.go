package controllers

import (
	"kwanjai/models"

	"github.com/gin-gonic/gin"
)

func VerifyEmail(ginContext *gin.Context) {
	verificationEmail := new(models.VerificationEmail)
	ginContext.ShouldBind(&verificationEmail)
	verificationEmail.UUID = ginContext.Param("UUID")
	status, message := verificationEmail.Verify()
	ginContext.JSON(status, gin.H{"message": message})
}
