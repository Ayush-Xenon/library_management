package controllers

import (
	"fmt"
	"library_management/initializers"
	"library_management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Enroll godoc
// @Summary Enroll User in Library
// @Description Enroll a user in a library
// @Tags library
// @Accept  json
// @Produce  json
// @Param  Authorization header string true "Bearer token"
// @Param  enroll body  models.EnrollRequest true  "enrollment data"
// @Success 200 {object} models.EnrollResponse "Enrollment successful"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Secrity BearerAuth
// @Router /auth/library/enroll [post]
func Enroll(c *gin.Context) {
	var enroll struct {
		LibraryID uint `binding:"required"`
	}
	if err := c.ShouldBindJSON(&enroll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "All fields required", "error": err.Error()})
		return
	}
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	userData := user.(models.User)
	fmt.Println(userData)
	var userLibrary models.UserLibraries
	initializers.DB.Model(models.UserLibraries{}).
		Where("user_id = ?", userData.ID).
		Where("library_id = ?", enroll.LibraryID).
		Find(&userLibrary)
	if userLibrary.UserID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already enrolled"})
		return
	}
	// var chk models.User
	// initializers.DB.Model(models.User{}).
	// 	Where("id = ?", userData.ID).
	// 	Where("role = ?", "user").Or("role = ?", "reader").
	// 	Find(&chk)
	// if chk.ID == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot enroll , user is owner or admin"})
	// 	return
	// }

	initializers.DB.Model(models.User{}).Where("id = ?", userData.ID).Update("Role", "reader")
	initializers.DB.Create(models.UserLibraries{UserID: userData.ID, LibraryID: enroll.LibraryID})
	c.JSON(http.StatusOK, gin.H{"data": "Enrolled"})
}
