package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"realTimeService/configuration"
	"strings"
)

func AuthMiddleware(cfg *configuration.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Here you can implement your authentication logic
		// For example, check for a valid token in the request header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		secretKey := []byte(cfg.SecretKey)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		sub, ok := claims["sub"].(string)
		if !ok {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("user_sub", sub)
		c.Set("auth_token", tokenString)
		// If the token is valid, proceed to the next handler
		c.Next()
	}
}
