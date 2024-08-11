package middlewares

import (
	"Todo/helpers"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			helpers.RespondWithError(c, http.StatusUnauthorized, "API token required", "401")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			helpers.RespondWithError(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}", "401")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			helpers.RespondWithError(c, http.StatusUnauthorized, err.Error(), "401")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			helpers.RespondWithError(c, http.StatusUnauthorized, "Invalid token claims", "401")
			c.Abort()
			return
		}

		// Extract user_id from the claims
		userID, ok := claims["user_id"] // JWT claims typically use float64 for numbers
		if !ok {
			helpers.RespondWithError(c, http.StatusUnauthorized, "Invalid token claims", "401")
			c.Abort()
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
