package controllers

import (
	"library_management/initializers"
	"library_management/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetBooksByPublisher(c *gin.Context) {
	publisher := c.Query("publisher")
	publisher = strings.ToUpper(publisher)
	var books []models.Book
	query := initializers.DB.Model(&models.Book{})
	if publisher != "" {
		query = query.Where("publisher ILIKE ?", "%"+publisher+"%")
	}
	query.Find(&books)

	if len(books) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": "No books found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": books})
}
