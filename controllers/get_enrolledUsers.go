package controllers

import (
	"library_management/initializers"
	"library_management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetEnrolledUsers(c *gin.Context) {

	id := c.Query("id")
	// title = strings.ToUpper(title)
	var usr []models.User
	// fmt.Println("dfcd", id)
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
	query := initializers.DB.Model(&models.User{}).
		Select("users.id", "users.name", "users.email", "users.role", "users.created_at", "users.updated_at", "users.contact_number").
		Joins("join user_libraries on user_libraries.user_id = users.id").
		Joins("join libraries on user_libraries.library_id = libraries.id").
		Where("libraries.id = ?", libUsr.LibraryID)
	if id != "" {
		query = query.Where("users.id =?", id)
	}
	query.Where("role=?", "reader").Find(&usr)

	if len(usr) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": "No User found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": usr})
}
