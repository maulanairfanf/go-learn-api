package middleware

import (
	"myapi/models"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTMiddlewareGin returns a Gin middleware for JWT authentication
func JWTMiddlewareGin() gin.HandlerFunc {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token format"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			c.Set("userID", claims.UserID)
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
			return
		}
	}
}
