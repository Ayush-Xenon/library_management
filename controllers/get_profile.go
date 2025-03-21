package controllers

import (
	"library_management/initializers"
	"library_management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {

	var pro models.User
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userData := user.(models.User)
	initializers.DB.Model(&models.User{}).
		Where("id=?", userData.ID).
		Find(&pro)
	if pro.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"data": "No libraries found"})
		return
	}
	// type libbb struct {
	// 	ID   uint
	// 	Name string
	// }
	// var libo []libbb
	// for _, l := range lib {
	// 	libo = append(libo, libbb{ID: l.ID, Name: l.Name})
	// }

	// initializers.DB.Model(models.User{}).Where("id = ?", userData.ID).Update("Role", "reader")
	// initializers.DB.Create(models.UserLibraries{UserID: userData.ID, LibraryID: enroll.LibraryID})
	c.JSON(http.StatusOK, gin.H{"data": pro})
}
