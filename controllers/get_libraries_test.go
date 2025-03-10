package controllers

import (
	"library_management/initializers"
	"library_management/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	initializers.DB = initializers.SetupTestDB()
	t.Cleanup(func() {
		initializers.CloseTestDB(initializers.DB)
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	return c, w
}

func TestGetLib(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func()
		expectedCode int
		expectedBody gin.H
	}{
		{
			name: "No libraries found",
			setupMocks: func() {},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"data": "No libraries found"},
		},
		{
			name: "Libraries found",
			setupMocks: func() {
				initializers.DB.Create(&models.Library{ID: 1, Name: "Library 1"})
				initializers.DB.Create(&models.Library{ID: 2, Name: "Library 2"})
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"data": []struct {
				ID   uint
				Name string
			}{
				{ID: 1, Name: "Library 1"},
				{ID: 2, Name: "Library 2"},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTest(t)
			tt.setupMocks()

			c.Request, _ = http.NewRequest(http.MethodGet, "/libraries", nil)
			GetLib(c)

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