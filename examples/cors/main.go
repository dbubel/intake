package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dbubel/intake/v2"
)

func main() {
	// Create a new Intake router
	router := intake.New()

	// Configure CORS middleware with default settings
	corsMiddleware := intake.CORS(intake.DefaultCORSConfig())

	// Add CORS middleware globally to all routes
	router.AddGlobalMiddleware(corsMiddleware)

	// Register a simple API endpoint
	router.AddEndpoint(http.MethodGet, "/api/data", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"message": "Hello from Intake!",
			"time":    time.Now().Format(time.RFC3339),
		}
		intake.RespondJSON(w, r, http.StatusOK, data)
	})

	// Register a POST endpoint to demonstrate different HTTP methods
	router.AddEndpoint(http.MethodPost, "/api/data", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"message": "Data received via POST",
			"time":    time.Now().Format(time.RFC3339),
		}
		intake.RespondJSON(w, r, http.StatusOK, data)
	})

	// Automatically add OPTIONS endpoints for all registered routes
	router.AddOptionsEndpoints()

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router.Mux,
	}

	fmt.Println("CORS-enabled server running at http://localhost:8080")
	fmt.Println("Try accessing the API from a browser or using fetch() in the console")
	fmt.Println("Example fetch: fetch('http://localhost:8080/api/data').then(r => r.json()).then(console.log)")

	// Start the server
	router.Run(server)
}
