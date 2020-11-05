package middlewares

import (
	"kwanjai/helpers"
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthorization middleware.
// Base authentication which always stores user object in Gin context.
// If token verification failed, anonymous user object is stored.
func JWTAuthorization() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		user := new(models.User)
		var token string
		extractedToken := strings.Split(ginContext.Request.Header.Get("Authorization"), "Bearer ")
		if len(extractedToken) != 2 {
			token = ""
		} else if extractedToken[0] != "Bearer " {
			token = ""
		} else {
			token = extractedToken[1]
		}
		passed, username, _, err := libraries.VerifyToken(token, "access")
		if !passed {
			user.MakeAnonymous()
			ginContext.Set("user", user)
			return
		}
		getUser, err := libraries.FirestoreFind("users", username)
		if err != nil {
			ginContext.AbortWithStatus(500)
			return
		}
		projects, _ := libraries.FirestoreSearch("projects", "Members", "array-contains", username)
		getUser.DataTo(user)
		user.Projects = len(projects)
		ginContext.Set("user", user)
	}
}

// AuthenticatedOnly disallows "anonymous" user.
func AuthenticatedOnly() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		if username == "anonymous" {
			ginContext.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
