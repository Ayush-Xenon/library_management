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

func Submit(c *gin.Context) {
	var issueId struct {
		ID uint
	}
	if err := c.BindJSON(&issueId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(issueId)

	var issue models.IssueRegistry
	initializers.DB.Model(&models.IssueRegistry{}).
		Where("id=?", issueId.ID).
		First(&issue)

	if issue.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Issue id not found"})
		return
	}

	if issue.IssueStatus != "lent" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book already returned"})
		return
	}

	var book models.Book
	initializers.DB.Model(&models.Book{}).
		Where("isbn=?", issue.ISBN).
		Where("lib_id=?", issue.LibId).
		First(&book)

	if book.LibID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
		return
	}

	issue.IssueStatus = "returned"
	issue.ReturnDate = time.Now()
	if err := initializers.DB.Save(&issue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	book.AvailableCopies += 1
	if err := initializers.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})

}
