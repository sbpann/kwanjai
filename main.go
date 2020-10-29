package main

import (
	"kwanjai/config"
	"kwanjai/controllers"
	"kwanjai/libraries"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	var err error

	config.BaseDirectory, err = os.Getwd()
	config.FrontendURL = "http://localhost:8000"
	config.EmailServicePassword, err = libraries.AccessSecretVersion("projects/978676563951/secrets/EmailServicePassword/versions/1")
	config.JWTAccessTokenSecretKey = "access"
	config.JWTRefreshTokenSecretKey = "refresh"
	config.JWTAccessTokenLifetime = time.Second * 4
	config.JWTRefreshTokenLifetime = time.Hour * 8
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/auth", controllers.AuthenticateTest)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
