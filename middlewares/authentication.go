package middlewares

import (
	"kwanjai/config"
	"kwanjai/libraries"
	"kwanjai/models"
	"log"

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
			if err != nil {
				log.Println(err.Error())
			}
			user.MakeAnonymous()
			ginContext.Set("user", user)
			ginContext.Next()
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
		ginContext.Next()
	}
}
