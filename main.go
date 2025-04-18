// filepath: library_management/main.go
package main

import (
	"library_management/controllers"
	_ "library_management/docs" // This is required to load the documentation files
	"library_management/initializers"
	"library_management/middlewares"
	"net/http"

	//	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Library Management API
// @version 1.0
// @description API for managing libraries and books, including user authentication and authorization.
// @host localhost:8081
// @BasePath /

func init() {
	initializers.LoadE()
	initializers.ConnectDB()
}

func main() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/signup", controllers.SignUp)

	r.POST("/login", controllers.Login)

	r.GET("/books", controllers.GetBooks)

	r.GET("/books/title", controllers.GetBooksByTitle)
	r.GET("/books/author", controllers.GetBooksByAuthor)
	r.GET("/books/publisher", controllers.GetBooksByPublisher)
	auth := r.Group("/auth")
	auth.Use(middlewares.CheckAuth())
	{
		auth.GET("/profile", controllers.GetProfile)
		auth.GET("/libraries", controllers.GetLib)
		auth.PATCH("/library/assign_admin", middlewares.CheckRole("owner"), controllers.AssignAdmin)
		auth.GET("/user", middlewares.CheckRole("owner"), controllers.GetUsers)
		auth.DELETE("/remove/user", controllers.RemoveUser)
		auth.DELETE("/delete/user", controllers.DeleteUser)
		auth.GET("/enrolleduser", controllers.GetEnrolledUsers)
		auth.POST("/library/create", middlewares.CheckRole("user"), controllers.CreateLibrary)
		auth.POST("/book/create", middlewares.CheckRole("admin"), controllers.CreateBook)
		auth.PATCH("/book/update", middlewares.CheckRole("admin"), controllers.UpdateBook)
		auth.POST("request/approve", middlewares.CheckRole("admin"), controllers.Approve)
		auth.PATCH("request/decline", middlewares.CheckRole("admin"), controllers.Decline)
		auth.PATCH("request/return", middlewares.CheckRole("admin"), controllers.Submit)
		auth.GET("/request/all", controllers.GetAllRequest)
		auth.GET("/issue/all", controllers.GetIssueReg)
		auth.POST("/library/enroll", middlewares.CheckRole2(), controllers.Enroll)
		auth.POST("request/raise", middlewares.CheckRole("reader"), controllers.RaiseRequest)
	}

	r.Run()
}
