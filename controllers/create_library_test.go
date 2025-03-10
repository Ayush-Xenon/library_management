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

// func setupTest(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
// 	gin.SetMode(gin.TestMode)
// 	initializers.DB = initializers.SetupTestDB()
// 	t.Cleanup(func() {
// 		initializers.CloseTestDB(initializers.DB)
// 	})

// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)

// 	return c, w
// }

func TestCreateLibrary(t *testing.T) {
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
		Role:          "user",
	}
	initializers.DB.Create(&testUser)
	tests := []struct {
		name         string
		body         models.LibraryInput
		setupMocks   func()
		currentUser  models.User
		expectedCode int
		expectedBody gin.H
	}{

		{
			name:         "Invalid JSON",
			body:         models.LibraryInput{},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "json: cannot unmarshal number into Go struct field LibraryInput.Name of type string"},
		},
		{
			name: "Library already exists",
			body: models.LibraryInput{
				Name: "TEST LIBRARYY",
			},
			setupMocks: func() {
				initializers.DB.Create(&models.Library{Name: "TEST LIBRARYY"})
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Library already exists"},
		},
		{
			name: "User not found",
			body: models.LibraryInput{
				Name: "TEST LIBRARY",
			},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "User not found"},
		},
		{
			name: "Successful creation",
			body: models.LibraryInput{
				Name: "TEST LIBRARY",
			},
			setupMocks: func() {},
			currentUser: models.User{
				ID:   1,
				Role: "user",
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"data": models.UserLibraries{
				UserID:    1,
				LibraryID: 1,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//c, w := setupTest(t)
			tt.setupMocks()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.body)
			c.Request, _ = http.NewRequest(http.MethodPost, "/libraries", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			if tt.currentUser.ID != 0 {
				c.Set("currentUser", tt.currentUser)
			}
			CreateLibrary(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			// var responseBody gin.H
			// json.Unmarshal(w.Body.Bytes(), &responseBody)
			// assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
	initializers.DB.Where("user_id=?", 1).Delete(&models.UserLibraries{})
	initializers.DB.Where("id=?", 1).Delete(&models.User{})
	initializers.DB.Where("name=?", "TEST LIBRARY").Delete(&models.Library{})
	initializers.DB.Where("name=?", "TEST LIBRARYY").Delete(&models.Library{})
	// initializers.DB.Where("name=?", "TEST LIBRARYY").Delete(&models.Library{})
}
