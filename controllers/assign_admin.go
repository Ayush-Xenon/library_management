package controllers

import (
	// 	//"fmt"
	"fmt"
	"library_management/initializers"
	"library_management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AssignAdmin godoc
// @Summary Assign Admin
// @Description Assign a user as an admin of a library
// @Tags library
// @Accept  json
// @Produce  json
// @Param  Authorization header string true "Bearer token"
// @Param  assignAdmin body  models.RequestID true  "assign admin data"
// @Success 200 {object} models.ErrorResponse "Admin assigned successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Security BearerAuth
// @Router /auth/library/assign_admin [patch]
func AssignAdmin(c *gin.Context) {
	var assign_admin struct {
		ID uint `binding:"required"`
	}
	if err := c.ShouldBindJSON(&assign_admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//c.JSON(http.StatusOK, gin.H{"data": assign_admin.ID})
	var userFound models.User
	initializers.DB.Model(models.User{}).
		Where("id=?", assign_admin.ID).Find(&userFound)
	if userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Owner not found"})
		return
	}
	userData := user.(models.User)
	fmt.Println(userData)
	var admins []struct {
		LibID int
	}

	var libUsr models.UserLibraries
	initializers.DB.Model(models.UserLibraries{}).
		Where("user_id = ?", userData.ID).
		Find(&libUsr)

	result := initializers.DB.Model(&models.Library{}).
		Select("libraries.id as lib_id").
		Joins("JOIN user_libraries ON user_libraries.library_id = libraries.id").
		Joins("JOIN users ON users.id = user_libraries.user_id").
		Where("user_libraries.library_id = ?", libUsr.LibraryID).
		Where("users.role = ?", "admin").
		Find(&admins)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	if len(admins) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Library already has an admin"})
		return
	}

	results := initializers.DB.Model(&models.User{}).
		Select("users.id as lib_id").
		Where("users.id = ?", assign_admin.ID).
		Where("users.role = ?", "user").
		Find(&admins)

	if results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": results.Error.Error()})
		return
	}
	if len(admins) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already is an admin or reader or owner"})
		return
	}

	initializers.DB.Create(&models.UserLibraries{UserID: assign_admin.ID, LibraryID: libUsr.LibraryID})
	initializers.DB.Model(models.User{}).Where("email = ?", userFound.Email).Update("Role", "admin")
	c.JSON(http.StatusOK, gin.H{"data": "admin assigned"})

}
