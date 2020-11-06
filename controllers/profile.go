package controllers

import (
	"kwanjai/helpers"
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProfilePicture() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		file, _, _ := ginContext.Request.FormFile("file")
		user, _ := ginContext.Get("user")
		if err := libraries.CloudStorageUpload(file, user.(*models.User).Username+".png"); err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if user.(*models.User).ProfilePicture == "https://storage.googleapis.com/kwanjai-a3803.appspot.com/anonymous.png" {
			libraries.FirestoreUpdateField("users", user.(*models.User).Username, "ProfilePicture", "https://storage.googleapis.com/kwanjai-a3803.appspot.com/"+user.(*models.User).Username+".png")
		}
		ginContext.JSON(http.StatusOK, gin.H{"message": "uploaded"})
	}
}

func MyProfile() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		user, _ := ginContext.Get("user") // user always exists
		ginContext.JSON(http.StatusOK, gin.H{
			"message": "Get profile successfully",
			"profile": user,
		})
	}
}

func UpdateProfile() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		profile := new(models.User)
		username := helpers.GetUsername(ginContext)
		ginContext.ShouldBindJSON(profile)
		if profile.Username != username {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		libraries.FirestoreUpdateFieldIfNotBlank("users", username, "Firstname", profile.Firstname)
		libraries.FirestoreUpdateFieldIfNotBlank("users", username, "Lastname", profile.Lastname)
		_, _, profile = profile.Finduser()
		ginContext.JSON(http.StatusOK, gin.H{
			"message": "Profile updated.",
		})
	}
}
