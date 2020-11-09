package helpers

import (
	"kwanjai/libraries"
	"kwanjai/models"

	"github.com/gin-gonic/gin"
)

// GetUsername fucntion returns username (string).
func GetUsername(ginContext *gin.Context) string {
	user, _ := ginContext.Get("user") // user always exists
	username := user.(*models.User).Username
	return username
}

// ExceedProjectLimit
func ExceedProjectLimit(ginContext *gin.Context) bool {
	user, _ := ginContext.Get("user") // user always exists
	plan := user.(*models.User).Plan
	projects := user.(*models.User).Projects
	switch plan {
	case "Starter":
		return !(projects < 1)
	case "Plus":
		return !(projects < 3)
	case "Pro":
		return false
	default:
		return true
	}
}

// IsSuperUser fucntion returns superuser status (bool).
func IsSuperUser(ginContext *gin.Context) bool {
	user, _ := ginContext.Get("user") // user always exists
	isSuperUser := user.(*models.User).IsSuperUser
	return isSuperUser
}

// IsProjectMember return membership status (bool) and error.
func IsProjectMember(username string, projectUUID string) bool {
	if projectUUID == "" {
		return false
	}
	project := new(models.Project)
	getProject, _ := libraries.FirestoreFind("projects", projectUUID) // projectUUID != "" ensures no error.
	getProject.DataTo(project)
	_, found := libraries.Find(project.Members, username)
	if !found {
		return false
	}
	return true
}

// IsOwner return ownership status (bool) and error.
func IsOwner(username string, objectType string, objectID string) bool {
	if objectID == "" || objectType == "" {
		return false
	}
	getObject, _ := libraries.FirestoreFind(objectType, objectID) // objectID != "", objectType != "" ensures no error.
	if !getObject.Exists() {
		return false
	}
	return getObject.Data()["User"].(string) == username
}
