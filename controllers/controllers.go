package controllers

import (
	"net/http"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// "gorm.io/driver/mysql"
	"golang.org/x/crypto/bcrypt"
    "fmt"
    "log"
    "govue/models"
    "strings"
    "time"
    "gorm.io/driver/postgres"
)


var db = make(map[string]string)

var DB *gorm.DB

var jwtKey = []byte("my_secret_key")
var tokens []string

type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func MiddleJWT(c *gin.Context) {
    bearerToken := c.Request.Header.Get("Authorization")
    reqToken := strings.Split(bearerToken, " ")[1]
    claims := &Claims{}
    tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            c.JSON(http.StatusUnauthorized, gin.H{
                "message": "unauthorized",
            })
            c.AbortWithStatus(400)
            return
        }
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "bad request",
        })
        c.AbortWithStatus(400)
        return
    }
    if !tkn.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{
            "message": "unauthorized",
        })
        c.AbortWithStatus(400)
        return
    }
    
    c.Next()
}

func GenerateJWT(expirationTime time.Time) (string, error) {
   
    claims := &Claims{
        Username: "username",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)



    return token.SignedString(jwtKey)

}

func Check(c *gin.Context) {
	password := []byte("passwords")

    // Hashing the password with the default cost of 10
    hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(hashedPassword))


    c.JSON(http.StatusOK, gin.H{
        "hash": string(hashedPassword),
    })  
}

func LoginDB(c *gin.Context){
	var user models.User
	var input models.User
    // Bind updated JSON data
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        c.AbortWithStatus(400)
        return
    }

    hashedPassword := []byte(user.Password)
    inputPassword := []byte(input.Password)
    fmt.Println(user.Password)
    fmt.Println(input.Password)

    // Comparing the password with the hash
    err := bcrypt.CompareHashAndPassword(hashedPassword, inputPassword)


    // Compare passwords (plain text, but should be hashed)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
        return
    }

    // Generate token
    token, err := AuthJWT(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token.Token,
        "expires_at": token.ExpiredAt,
    })  
}

func AuthJWT(user models.User) (models.PersonalToken, error) {
     expirationTime := time.Now().Add(5 * time.Minute)
     // Generate JWT token
    token, err := GenerateJWT(expirationTime)
    if err != nil {
        return models.PersonalToken{}, err
    }

    personalToken := models.PersonalToken{
        UserId:    user.ID,
        Token:     token,
        ExpiredAt: expirationTime,
    }

    // Save token to the database
    if err := DB.Create(&personalToken).Error; err != nil {
        return models.PersonalToken{}, err
    }

    return personalToken, nil
    
}

func InitDB() {
    dsn := "postgresql://postgres.bwionxlvltdwrtcrgfqp:Vue0895366740711@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres"
    var err error

    // ✅ Use postgres.Open instead of mysql.Open
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // // Auto migrate the Book model
    // err = DB.AutoMigrate(&models.Book{})
    // if err != nil {
    //     log.Fatal("Failed to auto-migrate:", err)
    // }
    // err1 := DB.AutoMigrate(&models.User{})
    // if err1 != nil {
    //     log.Fatal("Failed to auto-migrate:", err)
    // }
    // err2 := DB.AutoMigrate(&models.PersonalToken{})
    // if err2 != nil {
    //     log.Fatal("Failed to auto-migrate:", err)
    // }
}

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// func GetUser(c *gin.Context) {
// 	user := c.Params.ByName("name")
// 	value, ok := db[user]
// 	if ok {
// 		c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
// 	}
// }

// func AdminPost(c *gin.Context) {
// 	user := c.MustGet(gin.AuthUserKey).(string)

// 	// Parse JSON
// 	var json struct {
// 		Value string `json:"value" binding:"required"`
// 	}

// 	if c.Bind(&json) == nil {
// 		db[user] = json.Value
// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	}
// }

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
