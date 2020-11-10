package main

import (
	"context"
	_ "image/jpeg"
	_ "image/png"
	"kwanjai/config"
	"kwanjai/controllers"
	"kwanjai/libraries"
	"kwanjai/middlewares"
	"log"
	"net/http"
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
	if os.Getenv("PORT") != "" {
		config.Port = ":" + os.Getenv("PORT")
	} else {
		config.Port = ":8080"
	}
	libraries.InitializeGCP() // BaseDirectory need to be set before initialization.
	config.Context = context.Background()
	config.FrontendURL = "http://localhost:8080"
	config.BackendURL = "http://localhost:8080"
	config.FirebaseProjectID = "kwanjai-a3803"
	config.DefaultAuthenticationBackend = middlewares.JWTAuthorization()
	config.EmailServicePassword, err = libraries.AccessSecretVersion("projects/978676563951/secrets/EmailServicePassword/versions/1")
	config.EmailVerficationLifetime = time.Hour * 24 * 7
	config.JWTAccessTokenSecretKey, err = libraries.AccessSecretVersion("projects/978676563951/secrets/JWTAccessTokenSecretKey/versions/1")
	config.JWTRefreshTokenSecretKey, err = libraries.AccessSecretVersion("projects/978676563951/secrets/JWTRefreshTokenSecretKey/versions/1")
	config.JWTAccessTokenLifetime = time.Hour * 2
	config.JWTRefreshTokenLifetime = time.Hour * 8
	config.OmisePublicKey = "pkey_test_5lr43m5s2v9k5iohgnp"
	config.OmiseSecretKey = "skey_test_5lr43m5sdusunx0goly"
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
	api := ginEngine.Group("/api")
	authentication := api.Group("/authentication")
	authentication.POST("/login", controllers.Login())
	authentication.POST("/register", controllers.Register())
	authentication.POST("/logout", middlewares.AuthenticatedOnly(), controllers.Logout())
	authentication.POST("/verify_email/:UUID", controllers.VerifyEmail())
	authentication.POST("/resend_verification_email", controllers.ResendVerifyEmail())
	authentication.POST("/token/refresh", controllers.RefreshToken())
	authentication.GET("/token/verify", middlewares.AuthenticatedOnly(), controllers.TokenVerification())
	user := api.Group("/user")
	user.GET("/my_profile", middlewares.AuthenticatedOnly(), controllers.MyProfile())
	user.GET("/all", middlewares.AuthenticatedOnly(), controllers.AllUsernames())
	user.POST("/pay", middlewares.AuthenticatedOnly(), controllers.UpgradePlan())
	user.POST("/unsubscribe", middlewares.AuthenticatedOnly(), controllers.Unsubscribe())
	project := api.Group("/project")
	project.Use(middlewares.AuthenticatedOnly())
	{
		project.GET("/all", controllers.AllProject())
		project.POST("/new", controllers.NewProject())
		project.POST("/find", controllers.FindProject())
		project.PATCH("/update", controllers.UpdateProject())
		project.DELETE("/delete", controllers.DeleteProject())
	}
	board := api.Group("/board")
	board.Use(middlewares.AuthenticatedOnly())
	{
		board.POST("/all", controllers.AllBoard())
		board.POST("/new", controllers.NewBoard())
		board.POST("/find", controllers.FindBoard())
		board.PATCH("/update", controllers.UpdateBoard())
		board.DELETE("/delete", controllers.DeleteBoard())
	}
	post := api.Group("/post")
	post.Use(middlewares.AuthenticatedOnly())
	{
		post.POST("/all", controllers.AllPost())
		post.POST("/new", controllers.NewPost())
		post.PATCH("/update", controllers.UpdatePost())
		post.DELETE("/delete", controllers.DeletePost())
		post.POST("/comment/new", controllers.NewComment())
		post.PATCH("/comment/update", controllers.UpdateComment())
		post.DELETE("/comment/delete", controllers.DeleteComment())
	}
	ginEngine.Delims("$gin{", "}")
	ginEngine.LoadHTMLGlob("views/*")
	ginEngine.GET("/", func(ginContext *gin.Context) {
		ginContext.HTML(http.StatusOK, "index.html", gin.H{})
	})
	ginEngine.GET("/project/:ID", func(ginContext *gin.Context) {
		projectID := ginContext.Param("ID")
		ginContext.HTML(http.StatusOK, "project.html", gin.H{"projectId": projectID})
	})
	return ginEngine
}

func main() {
	setupServer()
	ginEngine := getServer(os.Getenv("GIN_MODE"))
	ginEngine.Run(config.Port)
}
