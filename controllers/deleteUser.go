package controllers

import (
	"fmt"
	"library_management/initializers"
	"library_management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteUser(c *gin.Context) {
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

	if pro.Role == "owner" && pro.Role == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot Perform This Action"})
		return
	}
	var issueReg []models.IssueRegistry
	initializers.DB.Model(&models.IssueRegistry{}).
		Where("reader_id=?", userData.ID).
		Where("issue_status=?", "lent").
		Find(&issueReg)

	if len(issueReg) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Has Some Issued Books"})
		return
	}

	initializers.DB.Model(&models.RequestEvent{}).
		Where("reader_id=?", userData.ID).
		Where("request_type=?", "required").
		Delete(&models.RequestEvent{})

	initializers.DB.Model(&models.UserLibraries{}).
		Where("user_id=?", userData.ID).
		Delete(&models.UserLibraries{})

	initializers.DB.Model(&models.User{}).
		Where("id=?", userData.ID).
		Delete(&models.User{})

	fmt.Println("fdnvkndakjv'nskjzdnszbckjbdzjvhchbbvj")
	c.JSON(http.StatusOK, gin.H{"data": "User Deleted Successfully"})
	return
}
