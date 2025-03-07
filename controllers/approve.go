package controllers

import(
	"library_management/initializers"
	"library_management/models"
	//"library_management/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Approve(c *gin.Context){
	var reqId struct{
		ID uint
	}
	if err:=c.BindJSON(&reqId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqReg models.RequestEvent
	initializers.DB.Model(&models.RequestEvent{}).
		Where("id=?",reqId.ID).
		First(&reqReg)

	if reqReg.ID==0{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request ID not found"})
		return
	}
	if reqReg.RequestType!="reqiured" {
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
		Where("isbn=?",reqReg.BookID).
		Where("lib_id=?",reqReg.LibID).
		First(&book)

	if book.LibID==0{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found in library"})
		return
	}

	if book.AvailableCopies<=0{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Currently not available"})
		//dfkljhgvrfvkrfvkrfjkgvrfnjhkgvnjhfrkvnjh

		//rfvjkrnvkrjnvkjnfrkjvgnjkfr
		return
	}
	book.AvailableCopies=book.AvailableCopies-1
	initializers.DB.Model(&models.Book{}).
	Where("isbn=?",reqReg.BookID).
	Where("lib_id=?",reqReg.LibID).
	Update("available_copies",book.AvailableCopies)

	initializers.DB.Model(&models.RequestEvent{}).
		Where("id=?",reqId.ID).
		Update("request_type","approved").
		Update("approver_id",userData.ID)

	var issue = models.IssueRegistry{
		ISBN               :reqReg.BookID,
		ReaderID           :reqReg.ReaderID,
		IssueApproverID    :userData.ID,
		IssueStatus        :"Lent",
		IssueDate          string
		ExpectedReturnDate string
		ReturnDate         string
		ReturnApproverID   uint
		LibId
	}


}