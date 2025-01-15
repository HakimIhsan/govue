package router

import (
	"github.com/gin-gonic/gin"

	"govue/controllers"
)

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	controllers.InitDB()
	r := gin.Default()
	r.Use(CORSMiddleware())
	// Get user value
	// r.GET("/user/:name", controllers.GetUser)
	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	r.POST("/login", controllers.LoginDB)

	authorized := r.Group("/", controllers.MiddleJWT)   

	// Ping test
	r.GET("/check", controllers.Check)
	authorized.GET("/ping", controllers.Ping)

	// CRUD Routes
    authorized.POST("/books", controllers.CreateBook)       // Create
    authorized.GET("/books", controllers.GetBooks)          // Read All
    authorized.GET("/books/:id", controllers.GetBookByID)   // Read One
    authorized.PUT("/books/:id", controllers.UpdateBook)    // Update
    authorized.DELETE("/books/:id", controllers.DeleteBook) // Delete

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	// authorized.POST("admin", controllers.AdminPost)

	return r
}

