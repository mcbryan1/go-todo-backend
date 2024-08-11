package router

import (
	"Todo/handlers"
	"Todo/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Apply CORS middleware
	applyCORS(router)

	// Define routes
	setupAuthRoutes(router)
	setupTodoRoutes(router)

	return router
}

// applyCORS configures CORS settings
func applyCORS(router *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Content-Type", "Authorization", "Accept", "Accept-Version", "Origin"}
	router.Use(cors.New(config))
}

// setupAuthRoutes configures authentication-related routes
func setupAuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/v1/auth")
	{
		authGroup.POST("/login", handlers.Login)
		authGroup.POST("/register", handlers.CreateUser)
	}
}

// setupTodoRoutes configures todo-related routes with authentication middleware
func setupTodoRoutes(router *gin.Engine) {
	todoGroup := router.Group("/v1/todo")
	todoGroup.Use(middlewares.AuthMiddleware())
	{
		todoGroup.GET("/get-todos", handlers.GetTodos)
		todoGroup.POST("/create-todo", handlers.CreateTodo)
		todoGroup.PUT("/edit-todo/:id", handlers.EditTodo)
		todoGroup.DELETE("delete-todo/:id", handlers.DeleteTodo)
	}
}
