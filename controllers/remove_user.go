package controllers

import (
	"fmt"
	"library_management/initializers"
	"library_management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RemoveUser(c *gin.Context) {
	var usrId = c.Query("id")

	var user models.User
	initializers.DB.Model((&models.User{})).
		Where("id=?", usrId).
		Find(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
	}
	var pro models.User
	c_user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userData := c_user.(models.User)
	initializers.DB.Model(&models.User{}).
		Where("id=?", userData.ID).
		Find(&pro)
	if pro.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"data": "User not found"})
		return
	}

	if user.Role == "admin" && pro.Role == "owner" {
		initializers.DB.Model(&models.UserLibraries{}).
			Where("user_id=?", user.ID).
			Delete(&models.UserLibraries{})

		initializers.DB.Model(&models.User{}).
			Where("id=?", user.ID).
			Update("role", "user")

		c.JSON(http.StatusOK, gin.H{"data": "Admin Removed Successfully"})
		return
	}

	var libUsr models.UserLibraries
	initializers.DB.Model(models.UserLibraries{}).
		Where("user_id = ?", userData.ID).
		Find(&libUsr)
	if user.Role == "reader" && pro.Role == "admin" {
		var issueReg []models.IssueRegistry
		initializers.DB.Model(&models.IssueRegistry{}).
			Where("reader_id=?", usrId).
			Where("lib_id=?", libUsr.LibraryID).
			Where("issue_status=?", "lent").
			Find(&issueReg)

		if len(issueReg) != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Has Some Issued Books"})
			return
		}

		initializers.DB.Model(&models.RequestEvent{}).
			Where("reader_id=?", usrId).
			Where("lib_id=?", libUsr.LibraryID).
			Where("request_type=?", "required").
			Update("request_type", "declined")

		initializers.DB.Model(&models.UserLibraries{}).
			Where("user_id=?", usrId).
			Where("library_id=?", libUsr.LibraryID).
			Delete(&models.UserLibraries{})

		var useLib []models.UserLibraries
		initializers.DB.Model(&models.UserLibraries{}).
			Where("user_id=?", usrId).
			Find(&useLib)

		if len(useLib) == 0 {
			initializers.DB.Model(&models.User{}).
				Where("id=?", usrId).
				Update("role", "user")
		}
		fmt.Println("fdnvkndakjv'nskjzdn")
		c.JSON(http.StatusOK, gin.H{"data": "Reader Removed Successfully"})
		return
	}
	// fmt.Println("rfjtrji")
	return
}
