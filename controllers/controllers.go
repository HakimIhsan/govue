package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
    "log"
    "govue/models"
)


var db = make(map[string]string)

var DB *gorm.DB

func InitDB() {
    dsn := "root:12345@tcp(127.0.0.1:8889)/govue"
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto migrate the Book model
    DB.AutoMigrate(&models.Book{})
}

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func GetUser(c *gin.Context) {
	user := c.Params.ByName("name")
	value, ok := db[user]
	if ok {
		c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
	}
}

func AdminPost(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)

	// Parse JSON
	var json struct {
		Value string `json:"value" binding:"required"`
	}

	if c.Bind(&json) == nil {
		db[user] = json.Value
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func Books(c *gin.Context) {
	// user := c.MustGet(gin.AuthUserKey).(string)

	c.JSON(http.StatusOK, gin.H{"value":c.Query("value")})

}

// CREATE a new book
func CreateBook(c *gin.Context) {
    var book models.Book

    // Bind JSON data to the models.Book struct
    if err := c.ShouldBindJSON(&book); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Insert data into the database
    if err := DB.Create(&book).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Book created successfully", "data": book})
}

func GetBooks(c *gin.Context) {
    var books []models.Book

    // Fetch all books from the database
    if err := DB.Find(&books).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return the list of books (no need for manual date parsing)
    c.JSON(http.StatusOK, books)
}


// READ a single book by ID
func GetBookByID(c *gin.Context) {
    id := c.Param("id")
    var book models.Book

    if err := DB.First(&book, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    c.JSON(http.StatusOK, book)
}

// UPDATE a book by ID
func UpdateBook(c *gin.Context) {
    id := c.Param("id")
    var book models.Book

    // Check if the book exists
    if err := DB.First(&book, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    var input models.Book

    // Bind updated JSON data
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Update the book
    DB.Model(&book).Updates(input)
    c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully", "data": book})
}

// DELETE a book by ID
func DeleteBook(c *gin.Context) {
    id := c.Param("id")
    var book models.Book

    // Check if the book exists
    if err := DB.First(&book, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    // Delete the book
    DB.Delete(&book)
    c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
