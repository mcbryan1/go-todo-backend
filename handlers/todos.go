package handlers

import (
	"Todo/helpers"
	"Todo/initializers"
	"Todo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTodos(c *gin.Context) {
	userID, ok, err := helpers.GetUserIDFromContext(c)
	if !ok || err != nil {
		if err == nil { // No specific error, assume unauthorized
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
	// Retrieve user_id from the context
	userID, exists := c.Get("user_id")
	if !exists {
		helpers.RespondWithError(c, http.StatusUnauthorized, "User not authenticated", "401")
		return
	}

	// Ensure userID is a UUID string and convert it to uuid.UUID
	typedUserIDStr, ok := userID.(string)
	if !ok {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve user ID from context", "500")
		return
	}

	// Parse request body
	var requestBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid request data", "001")
		return
	}

	// Create a new Todo instance
	newTodo := models.Todo{
		UserID:      typedUserIDStr,
		Title:       requestBody.Title,
		Description: requestBody.Description,
		Completed:   false,
	}

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
