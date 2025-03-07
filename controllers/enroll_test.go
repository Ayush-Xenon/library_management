package controllers

import (
	"bytes"
	"encoding/json"
	"library_management/initializers"
	"library_management/models"
	"net/http"
	"net/http/httptest"
	"testing"

	//"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // Importing SQLite driver for in-memory database
	"github.com/stretchr/testify/assert"
)

func TestEnroll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		body         gin.H
		currentUser  models.User
		setupMocks   func()
		expectedCode int
		expectedBody gin.H
	}{
		{
			name: "successful enrollment",
			body: gin.H{"LibraryID": 1},
			currentUser: models.User{
				ID:   1,
				Role: "user",
			},
			setupMocks: func() {
				initializers.DB = initializers.SetupTestDB()
				initializers.DB.Create(&models.User{ID: 1, Role: "user"})
				initializers.DB.Create(&models.Library{ID: 1})
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"data": "Enrolled"},
		},
		{
			name: "user not found",
			body: gin.H{"LibraryID": 1},
			setupMocks: func() {
				initializers.DB = initializers.SetupTestDB()
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "User not found"},
		},
		{
			name: "already enrolled",
			body: gin.H{"LibraryID": 1},
			currentUser: models.User{
				ID:   1,
				Role: "user",
			},
			setupMocks: func() {
				initializers.DB = initializers.SetupTestDB()
				initializers.DB.Create(&models.User{ID: 1, Role: "user"})
				initializers.DB.Create(&models.Library{ID: 1})
				initializers.DB.Create(&models.UserLibraries{UserID: 1, LibraryID: 1})
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Already enrolled"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.body)
			c.Request, _ = http.NewRequest(http.MethodPost, "/enroll", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			if tt.currentUser.ID != 0 {
				c.Set("currentUser", tt.currentUser)
			}

			Enroll(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			var responseBody gin.H
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
