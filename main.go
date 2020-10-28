package main

import (
	"kwanjai/config"
	"kwanjai/controllers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	var err error
	config.BaseDirectory, err = os.Getwd()
	config.FrontendURL = "http://localhost:8000"
	if err != nil {
		log.Println(err)
	}
	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
