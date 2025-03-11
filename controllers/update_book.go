package controllers

import (
	"library_management/initializers"
	"library_management/models"
	"library_management/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateBook godoc
// @Summary Update Book
// @Description Update the number of copies of a book in the library
// @Tags book
// @Accept  json
// @Produce  json
// @Param  Authorization header string true "Bearer token"
// @Param  updateBook body  models.UpdateBookRequest true  "update book data"
// @Success 200 {object} models.ErrorResponse "Book updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Security BearerAuth
// @Router /auth/book/update [patch]
func UpdateBook(c *gin.Context) {
	var update_book struct {
		ISBN   string `binding:"required"`
		Copies int    `binding:"required"`
	}

	if err := c.ShouldBindJSON(&update_book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var res models.ValidateOutput
	res = validators.ValidateISBN(update_book.ISBN)
	if !res.Result {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Message})
		return
	}
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userData := user.(models.User)
	var libUsr models.UserLibraries
	initializers.DB.Model(models.UserLibraries{}).
		Where("user_id = ?", userData.ID).
		Find(&libUsr)

	var book models.Book
	initializers.DB.Model(models.Book{}).
		Where("isbn=?", update_book.ISBN).
		Where("lib_id=?", libUsr.LibraryID).
		Find(&book)
	if book.ISBN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
		return
	}
	book.AvailableCopies = book.AvailableCopies + update_book.Copies
	book.TotalCopies = book.TotalCopies + update_book.Copies
	if update_book.Copies < 0 && book.TotalCopies < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "copies to be decreased must be less than total copies"})
		return
	}
	if update_book.Copies < 0 && book.AvailableCopies < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ccopies to be decreased must be less than available copies"})
		return
	}
	if book.TotalCopies == 0 {
		initializers.DB.
			Where("isbn=?", update_book.ISBN).
			Where("lib_id=?", libUsr.LibraryID).
			Delete(&models.Book{})
		c.JSON(http.StatusOK, gin.H{"msg": "Book Removed"})
		return
	}
	initializers.DB.Model(models.Book{}).
		Where("isbn=?", update_book.ISBN).
		Where("lib_id=?", libUsr.LibraryID).
		Update("total_copies", book.TotalCopies).
		Update("available_copies", book.AvailableCopies)
	c.JSON(http.StatusOK, gin.H{"data": "Book updated"})

}
