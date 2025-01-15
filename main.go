package handler

import (
	"govue/routes"
	"net/http"
)

// Exported handler for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	routers := router.SetupRouter()

	// Serve the request
	routers.ServeHTTP(w, r)
}

func main() {

	r := router.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
