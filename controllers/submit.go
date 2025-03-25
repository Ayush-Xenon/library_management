package controllers

import (
	"fmt"
	"time"

	//"library_management/controllers"
	"library_management/initializers"
	"library_management/models"

	//"time"

	//"library_management/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Submit godoc
// @Summary Submit Book
// @Description Submit a borrowed book
// @Tags book
// @Accept  json
// @Produce  json
// @Param  Authorization header string true "Bearer token"
// @Param  submit body  models.RequestID true  "submit data"
// @Success 200 {object} models.ErrorResponse "Book submitted successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Security BearerAuth
// @Router /auth/book/submit [post]
func Submit(c *gin.Context) {
	var issueId struct {
		ID uint `binding:"required"`
	}
	if err := c.BindJSON(&issueId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(issueId)

	var issue models.IssueRegistry
	initializers.DB.Model(&models.IssueRegistry{}).
		Where("id=?", issueId.ID).
		Find(&issue)

	if issue.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Issue id not found"})
		return
	}

	if issue.ReturnApproverID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book already returned"})
		return
	}

	var book models.Book
	initializers.DB.Model(&models.Book{}).
		Where("isbn=?", issue.ISBN).
		Where("lib_id=?", issue.LibId).
		Find(&book)

	if book.LibID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
		return
	}

	if err := initializers.DB.Model(&models.Book{}).
		Where("isbn=?", issue.ISBN).
		Where("lib_id=?", issue.LibId).
		Update("available_copies", book.AvailableCopies+1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userData := user.(models.User)
	if err := initializers.DB.Model(&models.IssueRegistry{}).
		Where("id=?", issueId.ID).
		Updates(map[string]interface{}{
			"issue_status":       "returned",
			"return_date":        time.Now(),
			"return_approver_id": userData.ID,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// initializers.DB.Model(&models.IssueRegistry{}).
	// 	Where("id=?", issueId.ID).
	// 	Update("return_approver_id", userData.ID)

	c.JSON(http.StatusOK, gin.H{"data": "Book returned successfully"})

}
