package controllers

import (
	"library_management/initializers"
	"library_management/models"
	"library_management/validators"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book in the library
// @Tags book
// @Accept  json
// @Produce  json
// @Param  Authorization header string true "Bearer token"
// @Param  book body  models.BookInput true  "book data"
// @Success 200 {object} models.ErrorResponse "Book created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Security BearerAuth
// @Router /auth/book/create [post]

func CreateBook(c *gin.Context) {
	var book models.BookInput
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "msg": "json: cannot unmarshal number into Go struct field BookInput.isbn of type string"})
		return
	}
	book.Authors = strings.ToUpper(book.Authors)
	book.Title = strings.ToUpper(book.Title)
	book.Publisher = strings.ToUpper(book.Publisher)

	var res models.ValidateOutput
	res = validators.ValidateISBN(book.ISBN)
	if !res.Result {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Message})
		return
	}

	var bookModel models.Book
	initializers.DB.Where("isbn = ?", book.ISBN).Find(&bookModel)
	if bookModel.ISBN != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book already exists"})
		return
	}

	if book.TotalCopies < 0 || book.AvailableCopies < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Available copies and Total copies cannot be less than 0"})
		return
	}
	if book.AvailableCopies != book.TotalCopies {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Available copies should be equal to  Total copies"})
		return
	}
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userData := user.(models.User)
	var libUsr models.UserLibraries
	initializers.DB.Model(&models.UserLibraries{}).
		Where("user_id = ?", userData.ID).
		Find(&libUsr)

	var newBook = models.Book{
		ISBN:            book.ISBN,
		LibID:           libUsr.LibraryID,
		Title:           book.Title,
		Authors:         book.Authors,
		Publisher:       book.Publisher,
		Version:         book.Version,
		TotalCopies:     book.TotalCopies,
		AvailableCopies: book.AvailableCopies,
	}

	initializers.DB.Create(&newBook)
	c.JSON(http.StatusOK, gin.H{"data": newBook})
}
