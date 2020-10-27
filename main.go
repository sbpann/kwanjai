package main

import (
	"gin-sandbox/controllers"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
