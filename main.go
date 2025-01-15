package handler

import (
	"govue/routes"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Exported handler for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	router := routes.SetupRouter()

	// Serve the request
	router.ServeHTTP(w, r)
}

func main() {

	r := router.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
