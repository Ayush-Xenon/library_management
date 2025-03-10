package controllers

import (
	"library_management/initializers"
	"library_management/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetBooksByAuthor(c *gin.Context) {
	author := c.Query("author")
	author = strings.ToUpper(author)
	var books []models.Book
	query := initializers.DB.Model(&models.Book{})
	if author != "" {
		query = query.Where("authors ILIKE ?", "%"+author+"%")
	}
	query.Find(&books)

	if len(books) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": "No books found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": books})
}
