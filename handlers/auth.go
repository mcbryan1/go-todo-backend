package handlers

import (
	"Todo/helpers"
	"Todo/initializers"
	"Todo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	_, user, tokenString, err := helpers.ProcessLogin(c)
	if err != nil {
		return
	}
	userResponse := helpers.CreateUserResponse(user)
	helpers.RespondWithSuccess(c, http.StatusOK, "Login Successful", "000", gin.H{
		"token": tokenString,
		"user":  userResponse,
	})
}

func CreateUser(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, "Invalid request body", "400")
		return
	}

	// Validate Request
	if err := helpers.ValidateRequest(req, "User"); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, err.Error(), "001")
		return
	}

	// Check if the user already exists
	email := req["email"].(string)
	if helpers.UserExistsByEmail(email) {
		helpers.RespondWithError(c, http.StatusConflict, "User already exists", "001")
		return
	}

	// Hash the password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to hash password", "500")
		return
	}

	// Create a new user
	newUser := models.User{
		Username: req["username"].(string),
		Email:    req["email"].(string),
		Password: string(hashPassword),
	}
	if err := initializers.DB.Create(&newUser).Error; err != nil {
		helpers.RespondWithError(c, http.StatusInternalServerError, "Failed to create user", "500")
	}

	helpers.RespondWithSuccess(c, http.StatusOK, "User created successfully", "000")

}
