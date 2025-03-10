package controllers

import (
	"library_management/initializers"
	"library_management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLib(c *gin.Context) {

	var lib []models.Library
	initializers.DB.Model(&models.Library{}).
		Find(&lib)
	if len(lib) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": "No libraries found"})
		return
	}
	type libbb struct {
		ID   uint
		Name string
	}
	var libo []libbb
	for _, l := range lib {
		libo = append(libo, libbb{ID: l.ID, Name: l.Name})
	}

	// initializers.DB.Model(models.User{}).Where("id = ?", userData.ID).Update("Role", "reader")
	// initializers.DB.Create(models.UserLibraries{UserID: userData.ID, LibraryID: enroll.LibraryID})
	c.JSON(http.StatusOK, gin.H{"data": libo})
}
