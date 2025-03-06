package controllers


import (
// 	//"fmt"
	"github.com/gin-gonic/gin"
	"library_management/initializers"
	"library_management/models"
	"net/http"
)

func AssignAdmin(c *gin.Context){
	 var assign_admin struct{
		id int
	 }
	if err := c.ShouldBindJSON(&assign_admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var userFound models.User
	initializers.DB.Where("email=?", userFound.Email).Find(&userFound)
	if userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	initializers.DB.Model(models.User{}).Where("email = ?", userFound.Email).Update("Role", "admin")
	c.JSON(http.StatusOK, gin.H{"data": "admin assigned"})

}