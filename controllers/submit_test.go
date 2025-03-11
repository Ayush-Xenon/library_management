package controllers

import (
	"bytes"
	"encoding/json"
	"library_management/initializers"
	"library_management/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type reqId struct {
	ID uint
}

func TestSubmit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	initializers.DB = initializers.SetupTestDB()
	defer initializers.CloseTestDB(initializers.DB)
	testIssue := models.IssueRegistry{
		ISBN:               "1234567890",
		ReaderID:           1,
		IssueApproverID:    2,
		IssueStatus:        "lent",
		ExpectedReturnDate: time.Now().AddDate(0, 0, 7),
		LibId:              1,
	}
	initializers.DB.Create(&testIssue)
	var iss models.IssueRegistry
	initializers.DB.Model(&models.IssueRegistry{}).
		Where("isbn=?", testIssue.ISBN).
		Find(&iss)

	var issId = iss.ID

	testIssue1 := models.IssueRegistry{
		ISBN:               "1234567891",
		ReaderID:           1,
		IssueApproverID:    2,
		IssueStatus:        "returned",
		ReturnApproverID:   2,
		ExpectedReturnDate: time.Now().AddDate(0, 0, 7),
		LibId:              1,
	}
	initializers.DB.Create(&testIssue1)
	var iss1 models.IssueRegistry
	initializers.DB.Model(&models.IssueRegistry{}).
		Where("isbn=?", testIssue1.ISBN).
		Find(&iss1)
	var issId1 = iss1.ID

	tests := []struct {
		name         string
		body         reqId
		setupMocks   func()
		currentUser  models.User
		expectedCode int
		expectedBody gin.H
	}{
		{
			name:         "Invalid JSON",
			body:         reqId{},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "json: cannot unmarshal string into Go struct field struct { ID uint } of type uint"},
		},
		{
			name: "Issue ID not found",
			body: reqId{
				ID: 1,
			},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Issue id not found"},
		},
		{
			name: "Book already returned",
			body: reqId{
				ID: issId1,
			},
			setupMocks: func() {
				initializers.DB.Create(&models.Book{ISBN: "1234567891", LibID: 1, AvailableCopies: 5, TotalCopies: 10})
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Book already returned"},
		},
		{
			name: "Book not found",
			body: reqId{
				ID: issId,
			},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Book not found"},
		},
		{
			name: "User not found",
			body: reqId{
				ID: issId,
			},
			setupMocks: func() {
				//initializers.DB.Create(&models.IssueRegistry{ID: 1, ISBN: "1234567890", LibId: 1, IssueStatus: "lent"})
				initializers.DB.Create(&models.Book{ISBN: "1234567890", LibID: 1, AvailableCopies: 5, TotalCopies: 10})
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"message": "User not found"},
		},
		{
			name: "Successful book return",
			body: reqId{
				ID: issId,
			},
			setupMocks: func() {
				//initializers.DB.Create(&models.IssueRegistry{ID: 1, ISBN: "1234567890", LibId: 1, IssueStatus: "lent"})
				//initializers.DB.Create(&models.Book{ISBN: "1234567892", LibID: 1, AvailableCopies: 5, TotalCopies: 10})
			},
			currentUser: models.User{
				ID: 2,
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"message": "Book returned successfully"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tt.currentUser.ID != 0 {
				c.Set("currentUser", tt.currentUser)
			}

			bodyBytes, _ := json.Marshal(tt.body)
			c.Request, _ = http.NewRequest(http.MethodPost, "/submit", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			Submit(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			// var responseBody gin.H
			// json.Unmarshal(w.Body.Bytes(), &responseBody)
			// assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
	initializers.DB.Where("isbn = ?", "1234567890").Delete(&models.Book{})
	initializers.DB.Where("isbn = ?", "1234567891").Delete(&models.Book{})
	initializers.DB.Where("reader_id=?", 1).Delete(&models.IssueRegistry{})
	//initializers.DB.Where("isbn=?", "1234567890").Delete(&models.IssueRegistry{})
}
