package middlewares

import (
	"kwanjai/config"
	"kwanjai/helpers"
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTAuthorization middleware.
// Base authentication which always stores user object in context.
// If token verification failed, anonymous user object is stored.
func JWTAuthorization() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		user := new(models.User)
		passed, username, _, err := libraries.VerifyToken(ginContext.Request.Header.Get("Authorization"), "access")
		if !passed {
			user.MakeAnonymous()
			ginContext.Set("user", user)
			return
		}
		firestoreClient, err := libraries.FirebaseApp().Firestore(config.Context)
		defer firestoreClient.Close()
		getUser, err := firestoreClient.Collection("users").Doc(username).Get(config.Context)
		if err != nil {
			ginContext.AbortWithStatus(500)
			return
		}
		getUser.DataTo(&user)
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
