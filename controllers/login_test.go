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

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	initializers.DB = initializers.SetupTestDB()
	defer initializers.CloseTestDB(initializers.DB)
	// Mock data
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("@Password123"), bcrypt.DefaultCost)
	user := models.User{
		ID:            1000,
		Name:          "TESTUSER",
		Password:      string(passwordHash),
		Email:         "test000@example.com",
		ContactNumber: "1234567890",
		Role:          "user",
	}
	initializers.DB.Create(&user)

	tests := []struct {
		name         string
		body         gin.H
		expectedCode int
		expectedBody gin.H
	}{
		{
			name: "Valid login",
			body: gin.H{
				"email":    "test000@example.com",
				"password": "@Password123",
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{
				// "token": func() string {
				// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				// 		"id":  user.ID,
				// 		"exp": time.Now().Add(time.Minute * 10).Unix(),
				// 	})
				// 	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
				// 	return tokenString
				// }(),
			},
		},
		{
			name: "Empty fields",
			body: gin.H{
				"email": "invalid000@example.com",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"msg": "All fields required",
			},
		},
		{
			name: "User not found",
			body: gin.H{
				"email":    "invalid@example.com",
				"password": "password123",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": "user not found",
			},
		},
		{
			name: "Invalid email",
			body: gin.H{
				"email":    "invalidexample.com",
				"password": "@Password123",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": "Invalid Email",
			},
		},
		{
			name: "Invalid password",
			body: gin.H{
				"email":    "test000@example.com",
				"password": "wrongpassword",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": "Password must contain at least 8 characters, a number, an uppercase letter, a lowercase letter and a special character",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.body)
			c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			Login(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			//assert.Contains(t, w.Body.String(), tt.expectedBody)
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)

			for key, value := range tt.expectedBody {
				assert.Equal(t, value, responseBody[key])
			}
		})
	}

	initializers.DB.Where("email=?", "test000@example.com").Delete(&models.User{})
	//initializers.DB.Where("email=?", 1).Delete(&models.User{})
}
