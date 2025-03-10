package controllers

import (
	"library_management/initializers"
	"library_management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {

	var books []models.Book
	initializers.DB.Model(&models.Book{}).
		Find(&books)
	if len(books) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": "No books found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": books})
}
