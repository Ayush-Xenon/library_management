package controllers

import (
	// "fmt"
	"library_management/initializers"
	"library_management/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetIssueReg(c *gin.Context) {
	types := c.Query("type")
	types = strings.ToLower(types)

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userData := user.(models.User)
	var usl models.UserLibraries
	initializers.DB.Model(&models.UserLibraries{}).
		Where("user_id=?", userData.ID).
		Find(&usl)

	var req []models.IssueRegistry
	if userData.Role == "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized access"})
		return
	}
	if userData.Role == "reader" {
		if types != "" {
			initializers.DB.Model(&models.IssueRegistry{}).
				Where("issue_status=?", types).
				Where("reader_id=?", usl.UserID).
				Find(&req)
			if len(req) == 0 {
				c.JSON(http.StatusOK, gin.H{"data": "No Issue Record Found"})
				return
			}
		} else {
			initializers.DB.Model(&models.IssueRegistry{}).
				Where("reader_id=?", usl.UserID).
				Find(&req)
			if len(req) == 0 {
				c.JSON(http.StatusOK, gin.H{"data": "No Issue Record Found"})
				return
			}
		}
	} else {
		if types != "" {
			initializers.DB.Model(&models.IssueRegistry{}).
				Where("issue_status=?", types).
				Where("lib_id=?", usl.LibraryID).
				Find(&req)
			if len(req) == 0 {
				c.JSON(http.StatusOK, gin.H{"data": "No Issue Record Found"})
				return
			}
		} else {
			initializers.DB.Model(&models.IssueRegistry{}).
				Where("lib_id=?", usl.LibraryID).
				Find(&req)
			if len(req) == 0 {
				c.JSON(http.StatusOK, gin.H{"data": "No Issue Record Found"})
				return
			}
		}
	}
	// fmt.Println(req[0].CreatedAt.Date())
	// for i := 0; i < len(req); i++ {
	// 	req[i].CreatedAt = req[i].CreatedAt.Local()
	// 	req[i].UpdatedAt = req[i].UpdatedAt.Local()
	// 	req[i].ReturnDate = req[i].ReturnDate.Local()
	// }
	c.JSON(http.StatusOK, gin.H{"data": req})
}
