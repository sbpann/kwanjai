package controllers

import (
	"kwanjai/config"
	"kwanjai/libraries"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AllUsernames() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		usernames := []string{}
		db := libraries.FirestoreDB()
		getUsername := db.Collection("users").Documents(config.Context)
		allUsernames, err := getUsername.GetAll()
		if err != nil {
			log.Panicln(err)
		}
		for _, user := range allUsernames {
			u := user.Data()["Username"].(string)
			usernames = append(usernames, u)
		}
		ginContext.JSON(http.StatusOK, gin.H{
			"message":   "Get username successfully",
			"usernames": usernames,
		})
	}
}
