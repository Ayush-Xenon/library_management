package main

import (
	// "fmt"
	// "log"
	// "net/http"
	// "os"
	"library_management/controllers"
	"library_management/initializers"

	// "library_management/models"
	"library_management/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadE()
	initializers.ConnectDB()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	auth := r.Group("/auth")
	auth.Use(middlewares.CheckAuth())
	{
		auth.POST("/library/create", controllers.CreateLibrary)
		//auth.GET("/library", controllers.GetLibrary)
		auth.PATCH("/library/assign_admin",middlewares.CheckRole("owner"),controllers.AssignAdmin)
	//r.POST("/library", middlewares.CheckAuth(), controllers.CreateLibrary)
	// r.GET("/library", controllers.GetLibrary)
	// r.POST("/book", controllers.CreateBook)
	// r.GET("/book", controllers.GetBook)
	}
	r.Run(":8081")
}
