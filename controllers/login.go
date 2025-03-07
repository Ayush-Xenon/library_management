package controllers

import (
	//"fmt"
	"net/http"
	"os"
	"strings"

	//"regexp"

	"library_management/initializers"
	"library_management/models"
	"library_management/validators"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "msg": "All fields required"})
		return
	}
	var res models.ValidateOutput
	res = validators.ValidateEmail(authInput.Email)
	if !res.Result {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Message})
		return
	}
	authInput.Email = strings.ToLower(authInput.Email)
	var userFound models.User
	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	res = validators.ValidatePassword(authInput.Password)
	if !res.Result {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Message})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}
	//expirationTime := time.Now().Add(1 * time.Minute)
	//c.SetCookie("token", token, int(expirationTime.Unix()), "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"token": token,
	})

}
