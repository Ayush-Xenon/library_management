package controllers

import (
	"library_management/initializers"
	"library_management/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetBooksByTitle(c *gin.Context) {
	title := c.Query("title")
	title = strings.ToUpper(title)
	var books []models.Book
	query := initializers.DB.Model(&models.Book{})
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	query.Find(&books)

	if len(books) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": "No books found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": books})
}
