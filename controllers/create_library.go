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

// CreateLibrary godoc
// @Summary Create a new library
// @Description Create a new library in the system
// @Tags library
// @Accept  json
// @Produce  json
// @Param  Authorization header string true "Bearer token"
// @Param  library body  models.LibraryInput true  "library data"
// @Success 200 {object} models.ErrorResponse "Library created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Security BearerAuth
// @Router /auth/library/create [post]

func CreateLibrary(c *gin.Context) {
	var library models.LibraryInput

	if err := c.ShouldBindJSON(&library); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	library.Name = strings.ToUpper(library.Name)
	var libraryModel models.Library
	initializers.DB.Where("name = ?", library.Name).Find(&libraryModel)
	if libraryModel.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Library Name Already Exists"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userData := user.(models.User)
	libraryModel.Name = library.Name
	initializers.DB.Create(&models.Library{Name: library.Name})

	var temp models.Library
	initializers.DB.Where("name = ?", library.Name).Find(&temp)

	initializers.DB.Create(&models.UserLibraries{UserID: userData.ID, LibraryID: temp.ID})

	initializers.DB.Model(models.User{}).Where("id = ?", userData.ID).Update("Role", "owner")

	// var owner = models.UserLibraries{
	// 	UserID:    userData.ID,
	// 	LibraryID: temp.ID,
	// }
	initializers.DB.Model(&models.User{}).
		Where("id=?", userData.ID).
		Find(&userData)
	c.Set("currentUser", userData)

	c.JSON(http.StatusOK, gin.H{"data": "Library Created Successfully"})

}
