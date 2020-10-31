package helpers

import (
	"kwanjai/models"

	"github.com/gin-gonic/gin"
)

// GetUsername fucntion returns username (string).
func GetUsername(ginContext *gin.Context) string {
	user, _ := ginContext.Get("user")
	username := user.(*models.User).Username
	return username
}

// IsSuperUser fucntion returns superuser staus of a user (bool).
func IsSuperUser(ginContext *gin.Context) bool {
	user, _ := ginContext.Get("user")
	isSuperUser := user.(*models.User).IsSuperUser
	return isSuperUser
}
