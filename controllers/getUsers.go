package controllers

import (
	// "fmt"
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
	// fmt.Println("dfcd", id)
	query := initializers.DB.Model(&models.User{})
	if id != "" {
		query = query.Where("id =?", id)
	}
	query.Where("role=?", "user").Find(&usr)

	if len(usr) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": "No User found"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userData := user.(models.User)
	if userData.Role == "user" || userData.Role == "reader" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized access"})
		return
	}
	var libUsr models.UserLibraries
	initializers.DB.Model(models.UserLibraries{}).
		Where("user_id = ?", userData.ID).
		Find(&libUsr)

	var admn []models.User
	// fmt.Println("dfcd", id)
	initializers.DB.Model(&models.User{}).
		Select("users.id", "users.name", "users.email", "users.role", "users.created_at", "users.updated_at", "users.contact_number").
		Joins("join user_libraries on user_libraries.user_id = users.id").
		Joins("join libraries on user_libraries.library_id = libraries.id").
		Where("libraries.id = ?", libUsr.LibraryID).
		Where("users.role=?", "admin").
		Find(&admn)

	// if len(usr) == 0 {
	// 	c.JSON(http.StatusOK, gin.H{"data": "No ",})
	// 	return
	// }
	fmt.Println("flsl;df", admn)

	c.JSON(http.StatusOK, gin.H{"data": usr, "admin": admn})
}
