package controllers

import (
	"library_management/initializers"
	"library_management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Enroll(c *gin.Context) {
	var enroll struct {
		LibraryID uint
	}
	if err := c.ShouldBindJSON(&enroll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userData := user.(models.User)
	var userLibrary models.UserLibraries
	initializers.DB.Model(models.UserLibraries{}).
		Where("user_id = ?", userData.ID).
		Where("library_id = ?", enroll.LibraryID).
		First(&userLibrary)
	if userLibrary.UserID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already enrolled"})
		return
	}
	// var chk models.User
	// initializers.DB.Model(models.User{}).
	// 	Where("id = ?", userData.ID).
	// 	Where("role = ?", "user").Or("role = ?", "reader").
	// 	First(&chk)
	// if chk.ID == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot enroll , user is owner or admin"})
	// 	return
	// }

	initializers.DB.Model(models.User{}).Where("id = ?", userData.ID).Update("Role", "reader")
	initializers.DB.Create(models.UserLibraries{UserID: userData.ID, LibraryID: enroll.LibraryID})
	c.JSON(http.StatusOK, gin.H{"data": "Enrolled"})
}
