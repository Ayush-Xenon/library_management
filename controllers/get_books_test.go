package controllers

import (
	//"encoding/json"
	"library_management/initializers"
	"library_management/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func TestGetBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	initializers.DB = initializers.SetupTestDB()
	defer initializers.CloseTestDB(initializers.DB)

	tests := []struct {
		name         string
		setupMocks   func()
		expectedCode int
		expectedBody gin.H
	}{
		{
			name: "No books found",
			setupMocks: func() {},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"data": []models.Book{}},
		},
		{
			name: "Books found",
			setupMocks: func() {
				initializers.DB.Create(&models.Book{ISBN: "1234567890", Title: "Book 1", Authors: "Author 1", Publisher: "Publisher 1", TotalCopies: 10, AvailableCopies: 10})
				initializers.DB.Create(&models.Book{ISBN: "0987654321", Title: "Book 2", Authors: "Author 2", Publisher: "Publisher 2", TotalCopies: 5, AvailableCopies: 5})
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"data": []models.Book{
				{ISBN: "1234567890", Title: "Book 1", Authors: "Author 1", Publisher: "Publisher 1", TotalCopies: 10, AvailableCopies: 10},
				{ISBN: "0987654321", Title: "Book 2", Authors: "Author 2", Publisher: "Publisher 2", TotalCopies: 5, AvailableCopies: 5},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/books", nil)
			GetBooks(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			// var responseBody gin.H
			// json.Unmarshal(w.Body.Bytes(), &responseBody)
			// assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
	initializers.DB.Exec("DELETE FROM books")
	initializers.DB.Exec("DELETE FROM request_events")
	initializers.DB.Exec("DELETE FROM user_libraries")
	initializers.DB.Exec("DELETE FROM users")
	initializers.DB.Exec("DELETE FROM libraries")
	initializers.DB.Exec("DELETE FROM issue_registries")
}