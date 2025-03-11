package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"library_management/initializers"
	"library_management/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUp(t *testing.T) {
	// Setup the router
	r := gin.Default()
	r.POST("/signup", SignUp)

	// Mock the database
	initializers.DB = initializers.SetupTestDB()
	defer initializers.CloseTestDB(initializers.DB)
	// Create a test user
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("@Password123"), bcrypt.DefaultCost)
	testUser := models.User{
		Name:          "TESTUSER",
		Password:      string(passwordHash),
		Email:         "test@example.com",
		ContactNumber: "1234567890",
		Role:          "user",
	}
	initializers.DB.Create(&testUser)

	// Define test cases
	tests := []struct {
		name         string
		input        models.AuthCreate
		expectedCode int
		expectedBody string
	}{
		{
			name: "Valid input",
			input: models.AuthCreate{
				Name:          "John Doe",
				Password:      "@Password123",
				Email:         "john@example.com",
				ContactNumber: "1234567890",
			},
			expectedCode: http.StatusCreated,
			expectedBody: `"SignUp successful"`,
		},
		{
			name: "Invalid email",
			input: models.AuthCreate{
				Name:          "John Doe",
				Password:      "@Password123",
				Email:         "invalid-email",
				ContactNumber: "1234567890",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `Invalid Email`,
		},
		{
			name: "Email already used",
			input: models.AuthCreate{
				Name:          "Test User",
				Password:      "@Password123",
				Email:         "test@example.com",
				ContactNumber: "1234567890",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"Email already used"`,
		},
		{
			name: "Invalid Password",
			input: models.AuthCreate{
				Name:          "Test User",
				Password:      "Password123",
				Email:         "test1@example.com",
				ContactNumber: "1234567890",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"Password must contain at least 8 characters, a number, an uppercase letter, a lowercase letter and a special character"`,
		},
		{
			name: "Invalid contact no",
			input: models.AuthCreate{
				Name:          "Test User",
				Password:      "@Password123",
				Email:         "test2@example.com",
				ContactNumber: "123456789",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"Invalid Contact No."`,
		},
		{
			name: "Invalid name",
			input: models.AuthCreate{
				Name:          "Test User ",
				Password:      "@Password123",
				Email:         "test3@example.com",
				ContactNumber: "1234567890",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"Invalid Name Format (3-50 characters, starts and ends with a letter and contains only letters and spaces and Find word should be at least 3 characters)"`,
		},
		{
			name: "Invalid Input",
			input: models.AuthCreate{
				Name:     "Test User ",
				Password: "@Password123",
				Email:    "test3@example.com",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"All fields required"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request body
			body, _ := json.Marshal(tt.input)

			// Create a request
			req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			w := httptest.NewRecorder()

			// Perform the request
			r.ServeHTTP(w, req)

			// Assert the response code
			assert.Equal(t, tt.expectedCode, w.Code)
			// Assert the response body
			// var response map[string]interface{}
			// json.Unmarshal(w.Body.Bytes(), &response)
			// assert.Equal(t, tt.expectedBody, response["message"])
			// Assert the response body
			// var response map[string]interface{}
			// json.Unmarshal(w.Body.Bytes(), &response)
			// //responseBody, _ := json.Marshal(response)
			// assert.Equal(t, response["message"], tt.expectedBody)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})

	}
	// Clean up the database
	// initializers.DB.Exec("DELETE FROM users")
	initializers.DB.Where("email=?", "test@example.com").Delete(&models.User{})
	initializers.DB.Where("email=?", "john@example.com").Delete(&models.User{})
}
