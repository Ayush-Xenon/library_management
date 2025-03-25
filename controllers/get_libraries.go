package controllers

import (
	"fmt"
	"library_management/initializers"
	"library_management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// type Dlib struct {
// 	models.Library
// 	stat string
// }

func GetLib(c *gin.Context) {

	var lib []models.Library
	initializers.DB.Model(&models.Library{}).
		Select("DISTINCT libraries.id, libraries.name").
		Joins("left JOIN user_libraries ON user_libraries.library_id = libraries.id").
		Where("libraries.id NOT IN (?)",
			initializers.DB.Model(&models.UserLibraries{}).
				Select("library_id").Where("user_id = ?", c.MustGet("currentUser").(models.User).ID)).
		Find(&lib)
	if len(lib) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": "No libraries found"})
		return
	}
	type libbb struct {
		ID   uint
		Name string
	}
	var libo []libbb
	for _, l := range lib {
		libo = append(libo, libbb{ID: l.ID, Name: l.Name})
	}
	fmt.Println(libo)
	// initializers.DB.Model(models.User{}).Where("id = ?", userData.ID).Update("Role", "reader")
	// initializers.DB.Create(models.UserLibraries{UserID: userData.ID, LibraryID: enroll.LibraryID})
	c.JSON(http.StatusOK, gin.H{"data": libo})
}
