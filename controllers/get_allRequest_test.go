package controllers

import (
	//"bytes"
	//"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"library_management/initializers"
	"library_management/models"
)


func TestGetAllRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	initializers.DB = initializers.SetupTestDB()
	defer initializers.CloseTestDB(initializers.DB)

	tests := []struct {
		name           string
		currentUser    models.User
		queryType      string
		expectedCode   int
		expectedBody   string
		setupMocks     func()
	}{
		{
			name:         "User not found",
			currentUser:  models.User{},
			queryType:    "",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"User not found"}`,
			setupMocks:   func() {},
		},
		{
			name:         "unauthorised",
			currentUser:  models.User{Role: "user", ID: 1},
			queryType:    "",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"data":"unauthorised"}`,
			setupMocks:   func() {},
		},
		{
			name:         "No books found",
			currentUser:  models.User{Role: "reader", ID: 1},
			queryType:    "",
			expectedCode: http.StatusOK,
			expectedBody: `{"data":"No books found"}`,
			setupMocks:   func() {},
		},
		{
			name:         "No books found with type user",
			currentUser:  models.User{Role: "reader", ID: 1},
			queryType:    "fiction",
			expectedCode: http.StatusOK,
			expectedBody: `{"data":"No books found"}`,
			setupMocks:   func() {},
		},
		{
			name:         "No books found with type",
			currentUser:  models.User{Role: "admin", ID: 2},
			queryType:    "fiction",
			expectedCode: http.StatusOK,
			expectedBody: `{"data":"No books found"}`,
			setupMocks:   func() {},
		},
		{
			name:         "No books found without type",
			currentUser:  models.User{Role: "admin", ID: 2},
			queryType:    "",
			expectedCode: http.StatusOK,
			expectedBody: `{"data":"No books found"}`,
			setupMocks:   func() {},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.currentUser.ID != 0 {
				c.Set("currentUser", tt.currentUser)
			}
			
			c.Request, _ = http.NewRequest(http.MethodGet, "/requests?type="+tt.queryType, nil)
			c.Request.Header.Set("Content-Type", "application/json")

			GetAllRequest(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			//assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}

	initializers.DB.Exec("DELETE FROM books")
	initializers.DB.Exec("DELETE FROM request_events")
	initializers.DB.Exec("DELETE FROM user_libraries")
	initializers.DB.Exec("DELETE FROM users")
	initializers.DB.Exec("DELETE FROM libraries")
	initializers.DB.Exec("DELETE FROM issue_registries")
}