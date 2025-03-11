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
)

func TestDecline(t *testing.T) {
	gin.SetMode(gin.TestMode)
	initializers.DB = initializers.SetupTestDB()
	defer initializers.CloseTestDB(initializers.DB)

	testReq := models.RequestEvent{
		BookID:      "1234567890",
		ReaderID:    1,
		RequestType: "required",
		LibID:       1,
	}
	initializers.DB.Create(&testReq)
	var reqs models.RequestEvent
	initializers.DB.Model(&models.RequestEvent{}).
		Where("id=?", testReq.ID).
		Find(&reqs)

	var reqsId = reqs.ID

	testReq1 := models.RequestEvent{
		BookID:      "1234567891",
		ReaderID:    1,
		RequestType: "approved",
		ApproverID:  2,
		LibID:       1,
	}
	initializers.DB.Create(&testReq1)
	var reqs1 models.RequestEvent
	initializers.DB.Model(&models.RequestEvent{}).
		Where("id=?", testReq1.ID).
		Find(&reqs1)

	var reqsId1 = reqs1.ID

	// testReq2 := models.RequestEvent{
	// 	BookID:      "1234567892",
	// 	ReaderID:    1,
	// 	RequestType: "required",
	// 	LibID:       1,
	// }
	// initializers.DB.Create(&testReq2)
	// var reqs2 models.RequestEvent
	// initializers.DB.Model(&models.RequestEvent{}).
	// 	Where("id=?", testReq2.ID).
	// 	Find(&reqs2)

	// var reqsId2 = reqs2.ID

	tests := []struct {
		name         string
		body         reqId
		setupMocks   func()
		currentUser  models.User
		expectedCode int
		expectedBody gin.H
	}{
		{
			name:         "Invalid JSON",
			body:         reqId{},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "json: cannot unmarshal string into Go struct field struct { ID uint } of type uint"},
		},
		{
			name: "Request not found",
			body: reqId{
				ID: 1,
			},
			setupMocks:   func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "Request not found"},
		},
		{
			name: "Request already approved or declined",
			body: reqId{
				ID: reqsId1,
			},
			setupMocks: func() {
				//initializers.DB.Create(&models.RequestEvent{ID: 1, RequestType: "approved"})
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "request already approved or declined"},
		},
		{
			name: "User not found",
			body: reqId{
				ID: reqsId,
			},
			setupMocks: func() {
				//initializers.DB.Create(&models.RequestEvent{ID: 1, RequestType: "required"})
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: gin.H{"error": "User not found"},
		},
		{
			name: "Successful decline",
			body: reqId{
				ID: reqsId,
			},
			setupMocks: func() {
				//initializers.DB.Create(&models.RequestEvent{ID: 1, RequestType: "required"})
			},
			currentUser: models.User{
				ID:   2,
				Role: "admin",
			},
			expectedCode: http.StatusOK,
			expectedBody: gin.H{"msg": "Request declined"},
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
			c.Request, _ = http.NewRequest(http.MethodPost, "/decline", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			Decline(c)

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
