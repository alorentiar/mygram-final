package main

import (
	"encoding/base64"
	"errors"
	"finalproject/core"
	"finalproject/database"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func hashPassword(password string) ([]byte, error) {
	cost := bcrypt.DefaultCost // Adjust cost as needed (higher cost takes longer)
	return bcrypt.GenerateFromPassword([]byte(password), cost)
}

func validateUser(user core.User) error {
	// Check for required fields (email, password)
	if user.Email == "" || user.Password == "" {
		return errors.New("Email and password are required")
	}

	return nil // No errors found
}

func RegisterUser(c *gin.Context) {
	// 1. Parse request body (replace with your actual User struct)
	var user core.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 2. Validate user data (use a validation library or custom logic)
	if err := validateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Hash password (use a secure hashing algorithm like bcrypt)
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	// user.Password = hashedPassword // Update user struct with hashed password
	user.Password = base64.StdEncoding.EncodeToString(hashedPassword)

	// 4. Connect to the database (assuming you have a database.NewPostgres() function)
	postgres := c.MustGet("user").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 5. Check for existing user (optional, based on your use case)
	var existingUser core.User
	err = db.Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil && existingUser.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // Ignore record not found error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for existing user"})
		return
	}

	// 6. Create the user record in the database
	tx := db.Begin() // Start a transaction
	err = tx.Create(&user).Error
	if err != nil {
		tx.Rollback() // Rollback transaction on error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	tx.Commit() // Commit transaction on success

	// 7. Send successful registration response
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})

}

func LoginUser(c *gin.Context) {
	// 1. Parse request body (replace with your actual User struct)
	var credentials LoginCredentials
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 2. Validate user credentials (use database lookup)
	var user core.User
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB
	err := db.Where("email = ?", credentials.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		}
		return
	}

	// 3. Compare password (use secure hashing algorithm like bcrypt)
	// if err := comparePassword(user.Password, credentials.Password); err != nil {
	//     c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
	//     return
	// }

	// 4. Generate JWT token (replace with your secret key and claims)
	// secretKey := "your_secret_key" // Replace with a secure, long-lived secret key stored securely (e.g., environment variable)
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing JWT secret key"})
		return
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,                              // Replace with the relevant user identifier claim
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // Set token expiration time (e.g., 1 hour)
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 5. Send successful login response
	c.JSON(http.StatusOK, gin.H{"token": token})
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// func comparePassword(hashedPassword []byte, password string) error {
// 	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
// }

func UpdateUser(c *gin.Context) {
	// 1. Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}

	// 2. Parse request body (replace with your actual User struct)
	var updatedUserData UserUpdate
	if err := c.BindJSON(&updatedUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 3. Validate updated user data (optional)
	// ... (implement validation logic if needed)

	// 4. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 5. Find user by ID
	var user User
	err := db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		}
		return
	}

	// 6. Update user data (consider using a struct with only updatable fields)
	user.Email = updatedUserData.Email // Update specific fields
	user.Name = updatedUserData.Name   // Update specific fields

	// 7. Save updated user in database
	err = db.Save(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// 8. Send successful update response
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

type UserUpdate struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	// ... other user fields
}

func DeleteUser(c *gin.Context) {
	// 1. Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}

	// 2. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 3. Find user by ID
	var user User
	err := db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		}
		return
	}

	// 4. Delete user from database
	err = db.Delete(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// 5. Send successful delete response
	c.JSON(http.StatusOK, gin.H{"message": "Success Delete"})
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
	// 1. Get photo ID from URL parameter
	photoID := c.Param("id")
	if photoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing photo ID"})
		return
	}

	// 2. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 3. Find photo by ID
	var photo Photo
	err := db.Where("id = ?", photoID).First(&photo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find photo"})
		}
		return
	}

	c.JSON(http.StatusOK, photo)
}

type Photo struct {
	ID          uint   `json:"id"`
	URL         string `json:"url"` // Adjust field name based on your storage strategy
	Description string `json:"description"`
}

func CreatePhoto(c *gin.Context) {
	// 1. Parse request body (replace with your actual Photo struct)
	var newPhoto Photo
	if err := c.BindJSON(&newPhoto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 2. (Optional) Validate photo data (e.g., URL or uploaded file)

	// 3. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 4. Save photo information in database
	err := db.Create(&newPhoto).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create photo"})
		return
	}

	// 5. (Optional) Handle image file upload and storage (replace based on your chosen storage strategy)
	//  - If storing image data in the database (not recommended):
	//    ... (logic to save image data to the database)

	//  - If storing image files separately (recommended):
	//    Use a library like `github.com/gin-gonic/gin/binding` to handle multipart form data and
	//    implement logic to save the uploaded image file to your chosen storage location.
	//    Update `newPhoto.URL` with the image file URL or path after successful storage.

	// 6. Send successful creation response
	c.JSON(http.StatusCreated, newPhoto)
}

func UpdatePhoto(c *gin.Context) {
	// 1. Get photo ID from URL parameter
	photoID := c.Param("id")
	if photoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing photo ID"})
		return
	}

	// 2. Parse request body (replace with your actual PhotoUpdate struct)
	var updatedPhotoData PhotoUpdate
	if err := c.BindJSON(&updatedPhotoData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 3. (Optional) Validate updated photo data

	// 4. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 5. Find photo by ID
	var photo Photo
	err := db.Where("id = ?", photoID).First(&photo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find photo"})
		}
		return
	}

	// 6. Update photo data (consider using a struct with only updatable fields)
	photo.Description = updatedPhotoData.Description // Update specific fields (replace with actual fields)

	// 7. (Optional) Handle image file update (replace based on your chosen storage strategy)
	//  - If storing image data in the database (not recommended):
	//    ... (logic to update image data in the database)

	//  - If storing image files separately (recommended):
	//    You might need additional logic to handle potential image file updates
	//    based on your chosen storage approach.

	// 8. Save updated photo in database
	err = db.Save(&photo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update photo"})
		return
	}

	// 9. Send successful update response
	c.JSON(http.StatusOK, gin.H{"message": "Photo updated successfully"})
}

type PhotoUpdate struct {
	Description string `json:"description"` // Replace with actual updatable fields
}

func DeletePhoto(c *gin.Context) {
	// 1. Get photo ID from URL parameter
	photoID := c.Param("id")
	if photoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing photo ID"})
		return
	}

	// 2. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 3. Find photo by ID
	var photo Photo
	err := db.Where("id = ?", photoID).First(&photo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find photo"})
		}
		return
	}

	// 4. Delete photo from database
	err = db.Delete(&photo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
		return
	}
}

// Comments
func GetAllComments(c *gin.Context) {
	// 1. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 2. (Optional) Apply filters based on query parameters (e.g., post ID, user ID)
	var comments []core.Comment
	query := db.Model(&core.Comment{}) // Start building the query

	// ... (implement logic to filter comments based on query parameters)

	// 3. Find all comments (or filtered comments)
	err := query.Find(&comments).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	// 4. Send successful response with all comments
	c.JSON(http.StatusOK, comments)
}

func GetOneComment(c *gin.Context) {
	// 1. Get comment ID from URL parameter
	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing comment ID"})
		return
	}

	// 2. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 3. Find comment by ID
	var comment core.Comment
	err := db.Where("id = ?", commentID).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find comment"})
		}
		return
	}

	// 4. Send successful response with the comment
	c.JSON(http.StatusOK, comment)
}

func CreateComment(c *gin.Context) {
	// 1. Parse request body (replace with your actual Comment struct)
	var newComment core.Comment
	if err := c.BindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 2. (Optional) Validate comment data

	// 3. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 4. Save comment in database
	err := db.Create(&newComment).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// 5. Send successful creation response
	c.JSON(http.StatusCreated, newComment)
}

func UpdateComment(c *gin.Context) {
	// 1. Get comment ID from URL parameter (continued)
	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing comment ID"})
		return
	}

	// 2. Parse request body (replace with your actual CommentUpdate struct)
	var updatedCommentData CommentUpdate
	if err := c.BindJSON(&updatedCommentData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 3. (Optional) Validate updated comment data

	// 4. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 5. Find comment by ID
	var comment core.Comment
	err := db.Where("id = ?", commentID).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find comment"})
		}
		return
	}

	// 6. Update comment data (consider using a struct with only updatable fields)
	// comment.Content = updatedCommentData.Content // Update specific fields (replace with actual fields)

	// 7. Save updated comment in database
	err = db.Save(&comment).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	// 8. Send successful update response
	c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
}

type CommentUpdate struct {
	Content string `json:"content"` // Replace with actual updatable fields
}

func DeleteComment(c *gin.Context) {
	// 1. Get comment ID from URL parameter
	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing comment ID"})
		return
	}

	// 2. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 3. Find comment by ID
	var comment core.Comment
	err := db.Where("id = ?", commentID).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find comment"})
		}
		return
	}

	// 4. Delete comment from database
	err = db.Delete(&comment).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	// 5. Send successful delete response
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// Social Media
func GetAllSocialMedia(c *gin.Context) {
	// 1. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 2. (Optional) Apply filters based on query parameters (e.g., user ID)
	var socialMediaData []core.SocialMedia // Replace with your actual model
	query := db.Model(&core.SocialMedia{}) // Start building the query

	// ... (implement logic to filter data based on query parameters)

	// 3. Find all social media data (or filtered data)
	err := query.Find(&socialMediaData).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get social media data"})
		return
	}

	// 4. Send successful response with all social media data
	c.JSON(http.StatusOK, socialMediaData)
}

func GetOneSocialMedia(c *gin.Context) {
	// 1. Get social media ID from URL parameter
	socialMediaID := c.Param("id")
	if socialMediaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing social media ID"})
		return
	}

	// 2. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 3. Find social media data by ID
	var socialMediaData core.SocialMedia // Replace with your actual model
	err := db.Where("id = ?", socialMediaID).First(&socialMediaData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Social media data not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get social media data"})
		}
		return
	}

	// 4. Send successful response with the social media data
	c.JSON(http.StatusOK, socialMediaData)
}

func CreateSocialMedia(c *gin.Context) {
	// 1. Parse request body (replace with your actual SocialMedia struct)
	var newSocialMediaData core.SocialMedia // Replace with your actual model
	if err := c.BindJSON(&newSocialMediaData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 2. (Optional) Validate social media data

	// 3. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 4. Save social media data in database
	err := db.Create(&newSocialMediaData).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create social media data"})
		return
	}

	// 5. Send successful creation response
	c.JSON(http.StatusCreated, newSocialMediaData)
}

func UpdateSocialMedia(c *gin.Context) {
	// 1. Get social media ID from URL parameter
	socialMediaID := c.Param("id")
	if socialMediaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing social media ID"})
		return
	}

	// 2. Parse request body (replace with your actual SocialMediaUpdate struct)
	var updatedSocialMediaData core.SocialMedia // Replace with struct containing updatable fields
	if err := c.BindJSON(&updatedSocialMediaData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 3. (Optional) Validate updated social media data

	// 4. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 5. Find social media data by ID
	var socialMediaData core.SocialMedia // Replace with your actual model
	err := db.Where("id = ?", socialMediaID).First(&socialMediaData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Social media data not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get social media data"})
		}
		return
	}

	// 6. Update social media data (consider using a struct with only updatable fields)
	// socialMediaData.Content = updatedSocialMediaData.Content  // Update specific fields (replace with actual fields)

	// 7. Save updated social media data in database
	err = db.Save(&socialMediaData).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update social media data"})
		return
	}

	// 8. Send successful update response
	c.JSON(http.StatusOK, gin.H{"message": "Social media data updated successfully"})
}

func DeleteSocialMedia(c *gin.Context) {
	// 1. Get social media ID from URL parameter
	socialMediaID := c.Param("id")
	if socialMediaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing social media ID"})
		return
	}

	// 2. Connect to database (replace with your database connection logic)
	postgres := c.MustGet("postgres").(*database.Postgres) // Assuming you store the connection in context
	db := postgres.DB

	// 3. Find social media data by ID
	var socialMediaData core.SocialMedia // Replace with your actual model
	err := db.Where("id = ?", socialMediaID).First(&socialMediaData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Social media data not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get social media data"})
		}
		return
	}

	// 4. Delete social media data from database
	err = db.Delete(&socialMediaData).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete social media data"})
		return
	}

	// 5. Send successful delete response
	c.JSON(http.StatusOK, gin.H{"message": "Social media data deleted successfully"})
}
