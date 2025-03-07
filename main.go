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
		auth.POST("/library/create", middlewares.CheckRole("user"), controllers.CreateLibrary)
		//auth.GET("/library", controllers.GetLibrary)
		auth.PATCH("/library/assign_admin", middlewares.CheckRole("owner"), controllers.AssignAdmin)
		//r.POST("/library", middlewares.CheckRole("user"), controllers.CreateLibrary)
		auth.POST("/library/enroll", middlewares.CheckRole2(), controllers.Enroll)
		//auth.POST("/library/enroll", middlewares.CheckRole("reader"), controllers.Enroll)
		// r.GET("/library", controllers.GetLibrary)
		auth.POST("/book/create", middlewares.CheckRole("admin"), controllers.CreateBook)
		// r.GET("/book", controllers.GetBook)
		auth.PATCH("/book/update", middlewares.CheckRole("admin"), controllers.UpdateBook)
		auth.POST("book/raise",middlewares.CheckRole("reader"),controllers.RaiseRequest)
	}
	r.Run(":8081")
}
