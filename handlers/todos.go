package handlers

import (
	"Todo/helpers"
	"Todo/initializers"
	"Todo/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTodos(c *gin.Context) {
	userID, ok, err := helpers.GetUserIDFromContext(c)
	if !ok || err != nil {
		if err == nil {
			helpers.RespondWithError(c, http.StatusUnauthorized, "User not authenticated", "401")
		} else {
			helpers.RespondWithError(c, http.StatusInternalServerError, err.Error(), "500")
		}
		return
	}

	var todos []models.Todo
	result := initializers.DB.Where("user_id = ?", userID).Find(&todos)

	if result.Error != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, result.Error, "001")
		return
	}
	todoResponse := helpers.CreateTodoResponses(todos)
	helpers.RespondWithSuccess(c, http.StatusOK, "Todos retrieved successfully", "000", todoResponse)
}

func CreateTodo(c *gin.Context) {
	// Retrieve user_id from the context using the helper function
	userID, ok, err := helpers.GetUserIDFromContext(c)
	if !ok || err != nil {
		if err == nil { // No specific error, assume unauthorized
			helpers.RespondWithError(c, http.StatusUnauthorized, "User not authenticated", "401")
		} else {
			helpers.RespondWithError(c, http.StatusInternalServerError, err.Error(), "500")
		}
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid request body", "400")
		return
	}
	// Validate Request
	if err := helpers.ValidateRequest(req, "Todo"); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error(), "001")
		return
	}

	// Create a new Todo instance
	newTodo := models.Todo{
		UserID:      userID, // Use the userID obtained from the helper function
		Title:       req["title"].(string),
		Description: req["description"].(string),
		Completed:   false,
	}

	fmt.Println("#############################todos", newTodo)

	// Save the new todo to the database
	result := initializers.DB.Create(&newTodo)
	if result.Error != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, result.Error.Error(), "001")
		return
	}

	todoResponse := helpers.CreateTodoResponse(newTodo)

	// Respond with success
	helpers.RespondWithSuccess(c, http.StatusCreated, "Todo created successfully", "000", todoResponse)
}

func EditTodo(c *gin.Context) {
	// Retrieve user_id from the context using the helper function
	userID, ok, err := helpers.GetUserIDFromContext(c)
	if !ok || err != nil {
		if err == nil { // No specific error, assume unauthorized
			helpers.RespondWithError(c, http.StatusUnauthorized, "User not authenticated", "401")
		} else {
			helpers.RespondWithError(c, http.StatusInternalServerError, err.Error(), "500")
		}
		return
	}

	// Parse and validate request body
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid request body", "400")
		return
	}

	// Extract Todo ID from the URL parameter
	todoID := c.Param("id")
	if todoID == "" {
		helpers.RespondWithError(c, http.StatusBadRequest, "Todo ID is required", "001")
		return
	}

	// Find the todo item in the database
	var todo models.Todo
	result := initializers.DB.First(&todo, "id = ? AND user_id = ?", todoID, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			helpers.RespondWithError(c, http.StatusNotFound, "Todo not found or not authorized", "001")
		} else {
			helpers.RespondWithError(c, http.StatusInternalServerError, result.Error.Error(), "001")
		}
		return
	}

	// Update the Todo item fields
	if title, ok := req["title"].(string); ok {
		todo.Title = title
	}
	if description, ok := req["description"].(string); ok {
		todo.Description = description
	}
	if completed, ok := req["completed"].(bool); ok {
		todo.Completed = completed
	}

	// Save the updated todo to the database
	result = initializers.DB.Save(&todo)
	if result.Error != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, result.Error.Error(), "001")
		return
	}

	// Respond with success
	todoResponse := helpers.CreateTodoResponse(todo)
	helpers.RespondWithSuccess(c, http.StatusOK, "Todo updated successfully", "000", todoResponse)
}
