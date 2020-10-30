package main

import (
	"context"
	"kwanjai/config"
	"kwanjai/controllers"
	"kwanjai/libraries"
	"kwanjai/middlewares"
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
	config.DefaultAuthenticationBackend = middlewares.JWTAuthorization()
	config.EmailServicePassword, err = libraries.AccessSecretVersion("projects/978676563951/secrets/EmailServicePassword/versions/1")
	config.EmailVerficationLifetime = time.Hour * 24 * 7
	config.JWTAccessTokenSecretKey, err = libraries.AccessSecretVersion("projects/978676563951/secrets/JWTAccessTokenSecretKey/versions/1")
	config.JWTRefreshTokenSecretKey, err = libraries.AccessSecretVersion("projects/978676563951/secrets/JWTRefreshTokenSecretKey/versions/1")
	config.JWTAccessTokenLifetime = time.Hour * 4
	config.JWTRefreshTokenLifetime = time.Hour * 8
	if err != nil {
		log.Fatalln(err)
	}
	// debugging area
	// user := new(models.User)
	// user.Username = "panithi.nakkhruea@pm.me"
	// status, message, user := models.Finduser(user)
	// fmt.Println(status, message, user)
	// debugging area

	r := gin.Default()
	r.Use(config.DefaultAuthenticationBackend)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)
	r.POST("/verify_email/:UUID", controllers.VerifyEmail)
	r.POST("/resend_verification_email", controllers.ResendVerifyEmail)
	r.POST("/token/refresh", controllers.RefreshToken)
	r.Run()
}
