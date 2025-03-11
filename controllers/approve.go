package controllers

import (
	"fmt"
	//"library_management/controllers"
	"library_management/initializers"
	"library_management/models"
	"time"

	//"library_management/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Approve godoc
// @Summary Approve Request
// @Description Approve a request event
// @Tags request
// @Accept  json
// @Produce  json
// @Param  Authorization header string true "Bearer token"
// @Param  reqId body  models.RequestID true  "request ID"
// @Success 200 {object} models.ApproveResponse "Request approved successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Security BearerAuth
// @Router /auth/request/approve [post]
func Approve(c *gin.Context) {
	var reqId struct {
		ID uint `binding:"required"`
	}
	if err := c.BindJSON(&reqId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//fmt.Println(reqId)
	var reqReg models.RequestEvent
	initializers.DB.Model(&models.RequestEvent{}).
		Where("id=?", reqId.ID).
		Find(&reqReg)
	fmt.Println(reqReg)
	if reqReg.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request ID not found"})
		return
	}
	if reqReg.RequestType != "required" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already approved or declined"})
		return
	}
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userData := user.(models.User)

	var book models.Book
	initializers.DB.Model(&models.Book{}).
		Where("isbn=?", reqReg.BookID).
		Where("lib_id=?", reqReg.LibID).
		Find(&book)

	if book.LibID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found in library"})
		return
	}

	if book.AvailableCopies <= 0 {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Currently not available"})
		//dfkljhgvrfvkrfvkrfjkgvrfnjhkgvnjhfrkvnjh
		initializers.DB.Model(&models.RequestEvent{}).
			Where("id=?", reqId.ID).
			Update("request_type", "declined")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not avilable request declined"})
		return
	}
	book.AvailableCopies = book.AvailableCopies - 1
	initializers.DB.Model(&models.Book{}).
		Where("isbn=?", reqReg.BookID).
		Where("lib_id=?", reqReg.LibID).
		Update("available_copies", book.AvailableCopies)

	initializers.DB.Model(&models.RequestEvent{}).
		Where("id=?", reqId.ID).
		Update("request_type", "approved").
		Update("approver_id", userData.ID)

	var t = time.Now().AddDate(0, 0, 7)

	var issue = models.IssueRegistry{
		ISBN:               reqReg.BookID,
		ReaderID:           reqReg.ReaderID,
		IssueApproverID:    userData.ID,
		IssueStatus:        "lent",
		ExpectedReturnDate: t,
		LibId:              reqReg.LibID,
	}
	initializers.DB.Create(&issue)
	c.JSON(http.StatusOK, gin.H{"msg": "Issued"})

}
