package main

import (
	"finalproject/core"
	"finalproject/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// User endpoints
	router.POST("/register", RegisterUser)
	router.POST("/login", LoginUser)
	router.PUT("/users/:id", UpdateUser)
	router.DELETE("/users/:id", DeleteUser)

	// Photo endpoints
	router.GET("/photos", GetAllPhotos)
	router.GET("/photos/:id", GetOnePhoto)
	router.POST("/photos", CreatePhoto)
	router.PUT("/photos/:id", UpdatePhoto)
	router.DELETE("/photos/:id", DeletePhoto)

	// Comment endpoints
	router.GET("/comments", GetAllComments)
	router.GET("/comments/:id", GetOneComment)
	router.POST("/comments", CreateComment)
	router.PUT("/comments/:id", UpdateComment)
	router.DELETE("/comments/:id", DeleteComment)

	// Social Media endpoints
	router.GET("/social-media", GetAllSocialMedia)
	router.GET("/social-media/:id", GetOneSocialMedia)
	router.POST("/social-media", CreateSocialMedia)
	router.PUT("/social-media/:id", UpdateSocialMedia)
	router.DELETE("/social-media/:id", DeleteSocialMedia)

	router.Run(":8080") // Start server on port 8080

	fmt.Println("Running on 8080!")

}

// Endpoint implementations (replace placeholders with actual logic and error handling)
// ...

// RegisterUser, LoginUser, UpdateUser, DeleteUser (as discussed previously)

func RegisterUser(c *gin.Context) {

}

func LoginUser(c *gin.Context) {

}

func UpdateUser(c *gin.Context) {

}

func DeleteUser(c *gin.Context) {

}

// Import the package that contains the definition of the `Photo` type

func GetAllPhotos(c *gin.Context) {
	var photos []core.Photo // Use the `models.Photo` type instead of `Photo`

	// Access the database connection from your Postgres struct
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you've stored the connection in context
	db := postgres.DB                                      // Get the gorm.DB instance

	// Find all photos using GORM's Find method
	err := db.Find(&photos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the list of photos
	c.JSON(http.StatusOK, photos)
}

func GetOnePhoto(c *gin.Context) {

}

func CreatePhoto(c *gin.Context) {

}

func UpdatePhoto(c *gin.Context) {

}

func DeletePhoto(c *gin.Context) {

}

// Comments
func GetAllComments(c *gin.Context) {

}

func GetOneComment(c *gin.Context) {

}

func CreateComment(c *gin.Context) {

}

func UpdateComment(c *gin.Context) {

}

func DeleteComment(c *gin.Context) {

}

// Social Media
func GetAllSocialMedia(c *gin.Context) {

}

func GetOneSocialMedia(c *gin.Context) {

}

func CreateSocialMedia(c *gin.Context) {

}

func UpdateSocialMedia(c *gin.Context) {

}

func DeleteSocialMedia(c *gin.Context) {

}
