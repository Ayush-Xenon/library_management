package controllers

import (
	"bytes"
	"encoding/json"
	"library_management/initializers"
	"library_management/models"

	//"library_management/validators"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRaiseRequest(t *testing.T) {
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
		Role:          "reader",
	}
	initializers.DB.Create(&testUser)
	tests := []struct {
		name         string
		body         models.RequestInput
		setupMocks   func()
		currentUser  models.User
		expectedCode int
		expectedBody gin.H
	}{
		{
			name: "Invalid JSON",
			body: models.RequestInput{
				LibID: 12345},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "json: cannot unmarshal number into Go struct field RequestInput.BookID of type string"},
		},
		{
			name: "Invalid ISBN",
			body: models.RequestInput{
				BookID: "invalid_isbn",
				LibID:  1,
			},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Invalid ISBN"},
		},
		{
			name: "User not found",
			body: models.RequestInput{
				BookID: "1234567890",
				LibID:  1,
			},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "User not found"},
		},
		{
			name: "Library not found",
			body: models.RequestInput{
				BookID: "1234567890",
				LibID:  1,
			},
			setupMocks: func() {
				//initializers.DB.Create(&models.UserLibraries{UserID: 1, LibraryID: 1})
			},
			currentUser: models.User{
				ID:   1,
				Role: "user",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Library not found"},
		},
		{
			name: "User not enrolled in library",
			body: models.RequestInput{
				BookID: "1234567890",
				LibID:  1,
			},
			setupMocks: func() {
				initializers.DB.Create(&models.Library{ID: 1, Name: "TEST"})
			},
			currentUser: models.User{
				ID:   1,
				Role: "user",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "User not enrolled in library"},
		},
		{
			name: "Book not found in library",
			body: models.RequestInput{
				BookID: "1234567890",
				LibID:  1,
			},
			setupMocks: func() {
				initializers.DB.Create(&models.UserLibraries{UserID: 1, LibraryID: 1})
				// initializers.DB.Create(&models.Library{ID: 1, Name: "Test Library"})
			},
			currentUser: models.User{
				ID:   1,
				Role: "user",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Book not found in library"},
		},
		{
			name: "Successful request",
			body: models.RequestInput{
				BookID: "1234567890",
				LibID:  1,
			},
			setupMocks: func() {
				//initializers.DB.Create(&models.UserLibraries{UserID: 1, LibraryID: 1})
				//initializers.DB.Create(&models.Library{ID: 1, Name: "T"})
				initializers.DB.Create(&models.Book{ISBN: "1234567890", LibID: 1})
			},
			currentUser: models.User{
				ID:   1,
				Role: "user",
			},
			expectedCode: http.StatusAccepted,
			expectedBody: gin.H{"msg": "Request issue raised"},
		},
		{
			name: "Book is already requested",
			body: models.RequestInput{
				BookID: "1234567890",
				LibID:  1,
			},
			setupMocks: func() {
				//initializers.DB.Create(&models.UserLibraries{UserID: 1, LibraryID: 1})
				//initializers.DB.Create(&models.Library{ID: 1, Name: "T"})
				initializers.DB.Create(&models.Book{ISBN: "1234567890", LibID: 1})
			},
			currentUser: models.User{
				ID:   1,
				Role: "user",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Book is already requested"},
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
			c.Request, _ = http.NewRequest(http.MethodPost, "/requests", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			RaiseRequest(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			// var responseBody gin.H
			// json.Unmarshal(w.Body.Bytes(), &responseBody)
			// assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
	initializers.DB.Where("isbn = ?", "1234567890").Delete(&models.Book{})
	initializers.DB.Where("library_id=?", 1).Delete(&models.UserLibraries{})
	initializers.DB.Where("id=?", 1).Delete(&models.User{})
	initializers.DB.Where("id=?", 1).Delete(&models.Library{})
	initializers.DB.Where("reader_id=?", 1).Delete(&models.RequestEvent{})

}
