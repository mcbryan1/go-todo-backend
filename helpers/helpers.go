package helpers

import (
	"Todo/initializers"
	"Todo/models"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserExists(email string) bool {
	var user models.User
	result := initializers.DB.Where("email = ?", email).First(&user)
	return result.Error == nil
}

func GenerateJWTToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	fmt.Println(token)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ProcessLogin(c *gin.Context) (req map[string]interface{}, user models.User, tokenString string, err error) {
	req, err = ParseRequest(c)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request data", "001")
	}

	user, err = GetUser(req["email"].(string))
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Invalid email or password", "001")
		return
	}

	if err = CheckPassword(user, req["password"].(string)); err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Invalid email or password", "001")
		return
	}

	tokenString, err = GenerateJWTToken(user)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Could not generate token", "500")
		return
	}

	return
}

func ParseRequest(c *gin.Context) (map[string]interface{}, error) {
	var req map[string]interface{}
	err := c.ShouldBindJSON(&req)
	return req, err
}

func RespondWithError(c *gin.Context, code int, message interface{}, resCode string) {
	c.AbortWithStatusJSON(code, gin.H{"resp_desc": message, "resp_code": resCode})
}

func GetUserIDFromContext(c *gin.Context) (string, bool, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false, nil
	}

	typedUserID, ok := userID.(string)
	if !ok {
		return "", false, fmt.Errorf("failed to retrieve user ID from context")
	}

	return typedUserID, true, nil
}

func RespondWithSuccess(c *gin.Context, code int, message interface{}, respCode string, data ...interface{}) {
	response := struct {
		RespCode string      `json:"resp_code"`
		RespDesc interface{} `json:"resp_desc"`
		Data     interface{} `json:"data"`
	}{
		RespCode: respCode,
		RespDesc: message,
		Data:     nil,
	}
	if len(data) > 0 {
		response.Data = data[0]
	}
	c.JSON(code, response)
}

func GetUser(email string) (models.User, error) {
	var user models.User
	err := initializers.DB.Where("email = ?", email).First(&user).Error

	return user, err

}

func CheckPassword(user models.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func ValidateRequest(req map[string]interface{}, req_type string) error {
	var requiredFields []string

	if req_type == "User" {
		requiredFields = []string{"email", "password", "username"}
	} else if req_type == "Todo" {
		requiredFields = []string{"title", "description"}
	} else {
		return fmt.Errorf("invalid request type")
	}

	for _, field := range requiredFields {
		if _, ok := req[field]; !ok {
			return fmt.Errorf("%s is required", field)
		}
		// Trim whitespace from the field value if it's a string
		strVal, ok := req[field].(string)
		if ok {
			strVal = strings.TrimSpace(strVal)
			if strVal == "" {
				return fmt.Errorf("%s cannot be empty", field)
			}
			req[field] = strVal
		}
	}

	// Validate Email
	email, ok := req["email"].(string)
	if !ok || !IsEmailValid(email) {
		return fmt.Errorf("invalid email")
	}

	return nil
}
func IsEmailValid(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	match := regexp.MustCompile(emailRegex).MatchString
	return match(email)
}

func UserExistsByEmail(email string) bool {
	var user models.User
	result := initializers.DB.Where("email = ?", email).First(&user)
	return result.Error == nil
}
