package controllers

import (
	"library_management/initializers"
	// "library_management/models"
	// "library_management/validators"
	"library_management/models"
	"library_management/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RaiseRequest godoc
// @Summary Raise Issue Request
// @Description Raise a request to issue a book
// @Tags request
// @Accept  json
// @Produce  json
// @Param  Authorization header string true "Bearer token"
// @Param  request body  models.RequestInput true  "request data"
// @Success 201 {object} models.ErrorResponse "Request raised successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Security BearerAuth
// @Router /auth/request/raise [post]
func RaiseRequest(c *gin.Context) {
	var req models.RequestInput

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var res models.ValidateOutput
	res = validators.ValidateISBN(req.BookID)
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

	var lib models.Library
	initializers.DB.Model(&models.Library{}).
		Where("id=?", req.LibID).
		Find(&lib)

	if lib.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Library not found"})
		return
	}

	var usrLib models.UserLibraries
	initializers.DB.Model(&models.UserLibraries{}).
		Where("user_id=?", userData.ID).
		Where("library_id=?", req.LibID).
		Find(&usrLib)

	if usrLib.LibraryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not enrolled in library"})
		return
	}

	var book models.Book
	initializers.DB.Model(&models.Book{}).
		Where("isbn=?", req.BookID).
		Where("lib_id=?", req.LibID).
		Find(&book)

	if book.LibID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found in library"})
		return
	}

	var chk models.RequestEvent
	initializers.DB.Model(&models.RequestEvent{}).
		Where("book_id=?", req.BookID).
		Where("lib_id=?", req.LibID).
		Where("request_type=?", "required").
		Where("reader_id=?", userData.ID).
		Find(&chk)
	if chk.LibID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book is already requested"})
		return
	}

	var reqIssue = models.RequestEvent{
		BookID:      req.BookID,
		ReaderID:    userData.ID,
		RequestType: "required",
		LibID:       req.LibID,
	}
	initializers.DB.Create(&reqIssue)
	c.JSON(http.StatusAccepted, gin.H{"msg": "Request issue raised"})

}
