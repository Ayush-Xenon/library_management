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

type up struct {
	ISBN   string
	Copies int
}

func TestUpdateBook(t *testing.T) {
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
		body         up
		setupMocks   func()
		currentUser  models.User
		expectedCode int
		expectedBody gin.H
	}{
		{
			name: "Invalid JSON",
			body: up{
				ISBN: "12345",
			},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "json: cannot unmarshal number into Go struct field struct { ISBN string; Copies int } of type string"},
		},
		{
			name: "Invalid ISBN",
			body: up{
				ISBN:   "invalid_isbn",
				Copies: 5,
			},
			setupMocks: func() {

			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Invalid ISBN"},
		},
		{
			name: "User not found",
			body: up{
				ISBN:   "1234567890",
				Copies: 5,
			},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "User not found"},
		},
		{
			name: "Book not found",
			body: up{
				ISBN:   "1234567890",
				Copies: 5,
			},
			setupMocks: func() {
				initializers.DB.Create(&testLib)
				initializers.DB.Create(&testUL)
			},
			currentUser: models.User{
				ID:   1,
				Role: "admin",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Book not found"},
		},
		{
			name: "Decrease copies less than total copies",
			body: up{
				ISBN:   "1234567890",
				Copies: -10,
			},
			setupMocks: func() {
				initializers.DB.Create(&models.Book{ISBN: "1234567890", LibID: 1, TotalCopies: 5, AvailableCopies: 3})
			},
			currentUser: models.User{
				ID:   1,
				Role: "admin",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "copies to be decreased must be less than total copies"},
		},
		{
			name: "Decrease copies less than available copies",
			body: up{
				ISBN:   "1234567890",
				Copies: -4,
			},
			setupMocks: func() {},
			currentUser: models.User{
				ID:   1,
				Role: "admin",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Ccopies to be decreased must be less than available copies"},
		},
		{
			name: "Successful book update",
			body: up{
				ISBN:   "1234567890",
				Copies: 5,
			},
			setupMocks: func() {},
			currentUser: models.User{
				ID:   1,
				Role: "admin",
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"data": "Book updated"},
		},
		{
			name: "Remove book",
			body: up{
				ISBN:   "1234567891",
				Copies: -10,
			},
			setupMocks: func() {
				initializers.DB.Create(&models.Book{ISBN: "1234567891", LibID: 1, TotalCopies: 10, AvailableCopies: 10})
			},
			currentUser: models.User{
				ID:   1,
				Role: "user",
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"msg": "Book Removed"},
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
			c.Request, _ = http.NewRequest(http.MethodPost, "/update-book", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			UpdateBook(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			// var responseBody gin.H
			// json.Unmarshal(w.Body.Bytes(), &responseBody)
			// assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
	initializers.DB.Where("library_id=?", 1).Delete(&models.UserLibraries{})
	initializers.DB.Where("id=?", 1).Delete(&models.User{})
	initializers.DB.Where("id=?", 1).Delete(&models.Library{})
	initializers.DB.Where("isbn = ?", "1234567890").Delete(&models.Book{})
	initializers.DB.Where("isbn = ?", "1234567891").Delete(&models.Book{})
}
