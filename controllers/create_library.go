package controllers

import (
	//"fmt"
	"library_management/initializers"
	"library_management/models"
	"net/http"
	"strings"

	//"os/user"

	"github.com/gin-gonic/gin"
)

func CreateLibrary(c *gin.Context) {
	var library models.LibraryInput

	if err := c.ShouldBindJSON(&library); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	library.Name = strings.ToUpper(library.Name)
	var libraryModel models.Library
	initializers.DB.Where("name = ?", library.Name).First(&libraryModel)
	if libraryModel.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Library already exists"})
		return
	}

	libraryModel.Name = library.Name
	initializers.DB.Create(&models.Library{Name: library.Name})

	var temp models.Library
	initializers.DB.Where("name = ?", library.Name).First(&temp)

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userData := user.(models.User)

	initializers.DB.Create(&models.UserLibraries{UserID: userData.ID, LibraryID: temp.ID})

	initializers.DB.Model(models.User{}).Where("id = ?", userData.ID).Update("Role", "owner")

	var owner = models.UserLibraries{
		UserID:    userData.ID,
		LibraryID: temp.ID,
	}
	c.JSON(http.StatusOK, gin.H{"data": owner})

}
