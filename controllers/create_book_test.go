package controllers

import (
	"bytes"
	"encoding/json"
	"library_management/initializers"
	"library_management/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	initializers.DB = initializers.SetupTestDB()
	defer initializers.CloseTestDB(initializers.DB)
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("@Password123"), bcrypt.DefaultCost)
	testUser := models.User{
		ID:            1,
		Name:          "TESTUSER",
		Password:      string(passwordHash),
		Email:         "test@example.com",
		ContactNumber: "1234567890",
		Role:          "admin",
	}
	testLib := models.Library{
		ID:   1,
		Name: "TEST",
	}
	initializers.DB.Create(&testUser)
	testUL := models.UserLibraries{UserID: 1, LibraryID: 1}
	tests := []struct {
		name         string
		body         models.BookInput
		currentUser  models.User
		setupMocks   func()
		expectedCode int
		expectedBody gin.H
	}{
		{
			name: "Invalid JSON",
			body: models.BookInput{
				ISBN: "1234567890"},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "json: cannot unmarshal number into Go struct field BookInput.isbn of type string"},
		},
		{
			name: "Invalid ISBN",
			body: models.BookInput{
				ISBN:            "invalid_isbn",
				Title:           "Test Book",
				Authors:         "Test Author",
				Publisher:       "Test Publisher",
				Version:         "1",
				TotalCopies:     10,
				AvailableCopies: 10,
			},
			// Uncomment and implement the mock setup if necessary
			// setupMocks: func() {
			// 	validators.ValidateISBN = func(isbn string) models.ValidateOutput {
			// 		return models.ValidateOutput{Result: false, Message: "Invalid ISBN"}
			// 	}
			// },
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Invalid ISBN"},
		},
		{
			name: "Book already exists",
			body: models.BookInput{
				ISBN:            "1234567890",
				Title:           "Test Book",
				Authors:         "Test Author",
				Publisher:       "Test Publisher",
				Version:         "1",
				TotalCopies:     10,
				AvailableCopies: 10,
			},
			setupMocks: func() {
				initializers.DB.Create(&models.Book{ISBN: "1234567890"})
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Book already exists"},
		},
		{
			name: "Negative copies",
			body: models.BookInput{
				ISBN:            "1234567890",
				Title:           "Test Book",
				Authors:         "Test Author",
				Publisher:       "Test Publisher",
				Version:         "1",
				TotalCopies:     -1,
				AvailableCopies: -1,
			},
			setupMocks: func() {
				initializers.DB.Where("isbn = ?", "1234567890").Delete(&models.Book{})
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Available copies and Total copies cannot be less than 0"},
		},
		{
			name: "Available copies not equal to total copies",
			body: models.BookInput{
				ISBN:            "1234567890",
				Title:           "Test Book",
				Authors:         "Test Author",
				Publisher:       "Test Publisher",
				Version:         "1",
				TotalCopies:     10,
				AvailableCopies: 5,
			},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Available copies should be equal to Total copies"},
		},
		{
			name: "User not found",
			body: models.BookInput{
				ISBN:            "1234567890",
				Title:           "Test Book",
				Authors:         "Test Author",
				Publisher:       "Test Publisher",
				Version:         "1",
				TotalCopies:     10,
				AvailableCopies: 10,
			},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "User not found"},
		},
		{
			name: "Successful creation",
			body: models.BookInput{
				ISBN:            "1234567890",
				Title:           "Test Book",
				Authors:         "Test Author",
				Publisher:       "Test Publisher",
				Version:         "1",
				TotalCopies:     10,
				AvailableCopies: 10,
			},
			currentUser: models.User{
				ID:   1,
				Role: "admin",
			},
			setupMocks: func() {
				// Ensure that no book exists before the test
				initializers.DB.Where("isbn = ?", "1234567890").Delete(&models.Book{})
				initializers.DB.Create(&testLib)
				initializers.DB.Create(&testUL)
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"data": models.Book{
				ISBN:            "1234567890",
				LibID:           1,
				Title:           "TEST BOOK",
				Authors:         "TEST AUTHOR",
				Publisher:       "TEST PUBLISHER",
				Version:         "1",
				TotalCopies:     10,
				AvailableCopies: 10,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.body)
			c.Request, _ = http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
			if tt.currentUser.ID != 0 {
				c.Set("currentUser", tt.currentUser)
			}
			// You should use the correct handler function, ensure CreateBook is defined
			CreateBook(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			// var responseBody gin.H
			// json.Unmarshal(w.Body.Bytes(), &responseBody)
			// assert.Equal(t, tt.expectedBody, responseBody)
			// var responseBody gin.H
			// json.Unmarshal(w.Body.Bytes(), &responseBody)
			// for key, value := range tt.expectedBody {
			// 	assert.Equal(t, value, responseBody[key])
			// }
		})
	}

	// Clean up all records created during the test
	initializers.DB.Where("isbn = ?", "1234567890").Delete(&models.Book{})
	initializers.DB.Where("library_id=?", 1).Delete(&models.UserLibraries{})
	initializers.DB.Where("id=?", 1).Delete(&models.User{})
	initializers.DB.Where("id=?", 1).Delete(&models.Library{})

}
