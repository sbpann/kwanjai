package main

import (
	"context"
	"fmt"
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
	config.Context = context.Background()
	config.FrontendURL = "http://localhost:8080"
	config.EmailServicePassword, err = libraries.AccessSecretVersion("projects/978676563951/secrets/EmailServicePassword/versions/1")
	config.EmailVerficationLifetime = time.Hour * 24 * 7
	config.JWTAccessTokenSecretKey, err = libraries.AccessSecretVersion("projects/978676563951/secrets/JWTAccessTokenSecretKey/versions/1")
	config.JWTRefreshTokenSecretKey, err = libraries.AccessSecretVersion("projects/978676563951/secrets/JWTRefreshTokenSecretKey/versions/1")
	config.JWTAccessTokenLifetime = time.Hour * 4
	config.JWTRefreshTokenLifetime = time.Hour * 8
	fmt.Println(config.JWTAccessTokenSecretKey, config.JWTRefreshTokenSecretKey)
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)
	r.POST("/verify/:UUID", controllers.VerifyEmail)
	r.POST("/token/refresh", controllers.RefreshToken)
	r.GET("/auth", controllers.AuthenticateTest)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
