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

func setupServer() {
	var err error
	if os.Getenv("GIN_MODE") == "" {
		os.Setenv("GIN_MODE", "debug")
	}
	config.BaseDirectory, err = os.Getwd()
	libraries.InitializeGCP() // BaseDirectory need to be set before initialization.
	config.Context = context.Background()
	config.FrontendURL = "http://localhost:8080"
	config.FirebaseProjectID = "kwanjai-a3803"
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
}

func getServer(mode string) *gin.Engine {
	if mode == "debug" {
		log.Println("running in debug mode.")
	} else if mode == "test" {
		gin.SetMode(gin.TestMode)
		log.Println("running in test mode.")
	} else if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
		log.Println("running in production mode.")
	}
	ginEngine := gin.Default()
	ginEngine.Use(config.DefaultAuthenticationBackend)
	ginEngine.POST("/login", controllers.Login())
	ginEngine.POST("/register", controllers.Register())
	ginEngine.POST("/logout", controllers.Logout())
	ginEngine.POST("/verify_email/:UUID", controllers.VerifyEmail())
	ginEngine.POST("/resend_verification_email", controllers.ResendVerifyEmail())
	ginEngine.POST("/token/refresh", controllers.RefreshToken())
	board := ginEngine.Group("/board")
	board.Use(middlewares.AuthenticatedOnly())
	{
		board.POST("/new", controllers.NewBoard())
		board.POST("/find", controllers.FindBoard())
		board.PATCH("/update", controllers.UpdateBoard())
		board.DELETE("/delete", controllers.DeleteBoard())
	}
	return ginEngine
}

func main() {
	setupServer()
	ginEngine := getServer(os.Getenv("GIN_MODE"))
	ginEngine.Run()
}
