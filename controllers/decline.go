package controllers

import (
	"fmt"
	// "library_management/controllers"
	// "library_management/initializers"
	"library_management/initializers"
	"library_management/models"

	// "time"

	// //"library_management/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Decline(c *gin.Context) {
	var reqId struct {
		ID uint
	}
	if err := c.BindJSON(&reqId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(reqId)

	var req models.RequestEvent
	initializers.DB.Model(&models.RequestEvent{}).
		Where("id=?", reqId.ID).
		First(&req)

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request not found"})
		return
	}

	if req.RequestType != "required" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request already approved or declined"})
		return
	}
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userData := user.(models.User)

	initializers.DB.Model(&models.RequestEvent{}).
		Where("id=?", req.ID).
		Update("request_type", "declined").
		Update("approver_id", userData.ID)

	c.JSON(http.StatusOK, gin.H{"msg": "Request declined"})

}
