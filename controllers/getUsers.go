package controllers

import (
	"fmt"
	"library_management/initializers"
	"library_management/models"
	"net/http"

	// "strings"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	id := c.Query("id")
	// title = strings.ToUpper(title)
	var usr []models.User
	fmt.Println("dfcd", id)
	query := initializers.DB.Model(&models.User{})
	if id != "" {
		query = query.Where("id =?", id)
	}
	query.Where("role=?", "user").Find(&usr)

	if len(usr) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": "No User found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": usr})
}
