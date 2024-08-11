package router

import (
	"Todo/handlers"
	"Todo/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Enable CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Content-Type", "Authorization", "Accept", "Accept-Version", "Origin"}
	router.Use(cors.New(config))

	// Setup Routes
	router.POST("/v1/auth/login", handlers.Login)
	router.POST("/v1/auth/register", handlers.CreateUser)
	router.GET("/v1/todo/get-todos", middlewares.AuthMiddleware(), handlers.GetTodos)
	router.POST("/v1/todo/create-todo", middlewares.AuthMiddleware(), handlers.CreateTodo)

	return router
}
