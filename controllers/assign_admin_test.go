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

type inp struct {
	ID uint
}

func TestAssignAdmin(t *testing.T) {
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
	initializers.DB.Create(&testUser)
	testUser1 := models.User{
		ID:            10,
		Name:          "TESTUSER",
		Password:      string(passwordHash),
		Email:         "testowner@example.com",
		ContactNumber: "1234567890",
		Role:          "owner",
	}
	initializers.DB.Create(&testUser1)
	initializers.DB.Create(&models.Library{ID: 10, Name: "TESTOWNER"})
	initializers.DB.Create(&models.UserLibraries{UserID: 10, LibraryID: 10})
	testUser2 := models.User{
		ID:            100,
		Name:          "TESTUSER",
		Password:      string(passwordHash),
		Email:         "testuser@example.com",
		ContactNumber: "1234567890",
		Role:          "user",
	}
	initializers.DB.Create(&testUser2)
	tests := []struct {
		name         string
		body         inp
		setupMocks   func()
		currentUser  models.User
		expectedCode int
		expectedBody gin.H
	}{
		{
			name:         "Invalid JSON",
			body:         inp{},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "json: cannot unmarshal string into Go struct field struct { ID uint } of type uint"},
		},
		{
			name: "User not found",
			body: inp{
				ID: 5,
			},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "user not found"},
		},
		{
			name: "Owner not found",
			body: inp{
				ID: 1,
			},
			setupMocks: func() {
				//initializers.DB.Create(&models.User{ID: 1, Email: "test@example.com"})
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Owner not found"},
		},
		{
			name: "User already is an admin or reader or owner",
			body: inp{
				ID: 1,
			},
			setupMocks: func() {
				//initializers.DB.Create(&models.User{ID: 1, Email: "test@example.com", Role: "user"})
				initializers.DB.Create(&models.Library{ID: 1, Name: "TEST"})
				initializers.DB.Create(&models.UserLibraries{UserID: 1, LibraryID: 1})
			},
			currentUser: models.User{
				ID:   testUser1.ID,
				Role: testUser1.Role,
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "User already is an admin or reader or owner"},
		},
		{
			name: "Successful admin assignment",
			body: inp{
				ID: testUser2.ID,
			},
			setupMocks: func() {
				//initializers.DB.Create(&models.User{ID: 1, Email: "test@example.com", Role: "user"})
				//initializers.DB.Create(&models.Library{ID: 1, Name: "Test Library"})
				//initializers.DB.Create(&models.UserLibraries{UserID: 2, LibraryID: 1})
			},
			currentUser: models.User{
				ID:   testUser1.ID,
				Role: testUser1.Role,
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"data": "admin assigned"},
		},
		{
			name: "Library already has an admin",
			body: inp{
				ID: testUser2.ID,
			},
			setupMocks: func() {
				//initializers.DB.Create(&models.User{ID: 1, Email: "test@example.com", Role: "user"})
				//	initializers.DB.Create(&models.Library{ID: 1, Name: "TEST"})
				//	initializers.DB.Create(&models.UserLibraries{UserID: 1, LibraryID: 1})
				//	initializers.DB.Create(&models.User{ID: 2, Email: "admin@example.com", Role: "admin"})
				//initializers.DB.Create(&models.UserLibraries{UserID: 2, LibraryID: 1})
			},
			currentUser: models.User{
				ID:   testUser1.ID,
				Role: testUser1.Role,
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Library already has an admin"},
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
			c.Request, _ = http.NewRequest(http.MethodPost, "/assign-admin", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			AssignAdmin(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			// var responseBody gin.H
			// json.Unmarshal(w.Body.Bytes(), &responseBody)
			// assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
	initializers.DB.Where("library_id=?", 1).Delete(&models.UserLibraries{})
	initializers.DB.Where("library_id=?", 10).Delete(&models.UserLibraries{})
	initializers.DB.Where("id=?", 1).Delete(&models.User{})
	initializers.DB.Where("id=?", 10).Delete(&models.User{})
	initializers.DB.Where("id=?", 100).Delete(&models.User{})
	initializers.DB.Where("id=?", 1).Delete(&models.Library{})
	initializers.DB.Where("id=?", 10).Delete(&models.Library{})
	//	initializers.DB.Where("book_id=?", "1234567890").Delete(&models.RequestEvent{})
}
