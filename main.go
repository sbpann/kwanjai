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
	config.BackendURL = "http://localhost:8080"
	config.FirebaseProjectID = "kwanjai-a3803"
	config.DefaultAuthenticationBackend = middlewares.JWTAuthorization()
	config.EmailServicePassword, err = libraries.AccessSecretVersion("projects/978676563951/secrets/EmailServicePassword/versions/1")
	config.EmailVerficationLifetime = time.Hour * 24 * 7
	config.JWTAccessTokenSecretKey, err = libraries.AccessSecretVersion("projects/978676563951/secrets/JWTAccessTokenSecretKey/versions/1")
	config.JWTRefreshTokenSecretKey, err = libraries.AccessSecretVersion("projects/978676563951/secrets/JWTRefreshTokenSecretKey/versions/1")
	config.JWTAccessTokenLifetime = time.Hour * 2
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
	api := ginEngine.Group("/api")
	// authentication := api.Group("/authentication") // uncomment this for built-in frontend mode
	authentication := ginEngine.Group("/authentication") // uncomment this for dedicated frontend mode
	authentication.POST("/login", controllers.Login())
	authentication.POST("/register", controllers.Register())
	authentication.POST("/logout", middlewares.AuthenticatedOnly(), controllers.Logout())
	authentication.POST("/verify_email/:UUID", controllers.VerifyEmail())
	authentication.POST("/resend_verification_email", controllers.ResendVerifyEmail())
	authentication.POST("/token/refresh", controllers.RefreshToken())
	authentication.GET("/token/verify", controllers.TokenVerification())
	user := api.Group("/user")
	user.GET("/my_profile", middlewares.AuthenticatedOnly(), controllers.MyProfile())
	user.GET("/all", middlewares.AuthenticatedOnly(), controllers.AllUsernames())
	// project := api.Group("/project") // uncomment this for built-in frontend mode
	project := ginEngine.Group("/project") // uncomment this for dedicated frontend mode
	project.Use(middlewares.AuthenticatedOnly())
	{
		project.GET("/all", controllers.AllProject())
		project.POST("/new", controllers.NewProject())
		project.POST("/find", controllers.FindProject())
		project.PATCH("/update", controllers.UpdateProject())
		project.DELETE("/delete", controllers.DeleteProject())
	}
	// board := api.Group("/board") // uncomment this for built-in frontend mode
	board := ginEngine.Group("/board") // uncomment this for dedicated frontend mode
	board.Use(middlewares.AuthenticatedOnly())
	{
		board.POST("/all", controllers.AllBoard())
		board.POST("/new", controllers.NewBoard())
		board.POST("/find", controllers.FindBoard())
		board.PATCH("/update", controllers.UpdateBoard())
		board.DELETE("/delete", controllers.DeleteBoard())
	}
	// post := api.Group("/post") // uncomment this for built-in frontend mode
	post := ginEngine.Group("/post") // uncomment this for dedicated frontend mode
	post.Use(middlewares.AuthenticatedOnly())
	{
		post.POST("/all", controllers.AllPost())
		post.POST("/new", controllers.NewPost())
		post.PATCH("/find", controllers.FindPost())
		post.DELETE("/delete", controllers.DeletePost())
		post.POST("/comment/new", controllers.NewComment())
		post.PATCH("/comment/update", controllers.UpdateComment())
		post.DELETE("/comment/delete", controllers.DeleteComment())
	}
	// uncomment these for built-in frontend mode
	// ginEngine.Delims("$gin{", "}")
	// ginEngine.LoadHTMLGlob("views/*")
	// ginEngine.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", gin.H{})
	// })
	// ginEngine.GET("/project/:ID", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "project.html", gin.H{"id": c.Param("ID")})
	// })
	// uncomment these for built-in frontend mode
	return ginEngine
}

func main() {
	setupServer()
	ginEngine := getServer(os.Getenv("GIN_MODE"))
	ginEngine.Run()
}
