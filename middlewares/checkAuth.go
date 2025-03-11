package middlewares

import (
	"fmt"
	"library_management/initializers"
	"library_management/models"

	//"task/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authToken := strings.Split(authHeader, " ")
		fmt.Println(authToken)
		if len(authToken) != 2 || authToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := authToken[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.User
		initializers.DB.Where("ID=?", claims["id"]).Find(&user)

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("currentUser", user)

		c.Next()
	}
}

func CheckRole(s string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("currentUser").(models.User)
		if user.Role != s {
			if user.Role == "user" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action , Find enroll or create"})
			}
			if user.Role == "admin" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action , admin cannot be a reader or owner"})
			}
			if user.Role == "owner" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action , owner cannot be a reader or admin"})
			}
			if user.Role == "reader" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action , reader cannot be owner or admin"})
			}
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}

}

func CheckRole2() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("currentUser").(models.User)

		// if user.Role == "user" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action , Find enroll or create"})
		// }
		if user.Role == "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action , admin cannot be a reader or owner"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if user.Role == "owner" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action , owner cannot be a reader or admin"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}

}
