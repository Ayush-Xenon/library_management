package controllers

import (
	//"fmt"
	"net/http"
	//	"os"

	//"regexp"

	"library_management/initializers"
	"library_management/models"
	"library_management/validators"

	//"time"
	"strings"

	"github.com/gin-gonic/gin"

	//	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {

	var authInput models.AuthCreate
	//fmt.Println("input rec")
	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "All fields required"})
		return
	}
	var response models.ValidateOutput
	response = validators.ValidateEmail(authInput.Email)
	if !response.Result {
		c.JSON(http.StatusBadRequest, gin.H{"message": response.Message})
		return
	}
	response = validators.ValidatePassword(authInput.Password)
	if !response.Result {
		c.JSON(http.StatusBadRequest, gin.H{"message": response.Message})
		return
	}
	response = validators.ValidateName(authInput.Name)
	if !response.Result {
		c.JSON(http.StatusBadRequest, gin.H{"message": response.Message})
		return
	}
	response = validators.ValidatePhone(authInput.ContactNumber)
	if !response.Result {
		c.JSON(http.StatusBadRequest, gin.H{"message": response.Message})
		return
	}
	authInput.Email = strings.ToLower(authInput.Email)
	authInput.Name = strings.ToUpper(authInput.Name)

	var userFound models.User
	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email already used"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := models.User{
		Name:          authInput.Name,
		Password:      string(passwordHash),
		Email:         authInput.Email,
		ContactNumber: authInput.ContactNumber,
		Role:          "user",
	}

	initializers.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user, "message": "SignUp successful"})

}
