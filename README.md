# Intake

Intake is a lightweight, flexible HTTP router for Go applications with middleware support. It provides a simple API for defining routes, handling requests, and applying middleware.

## Features

- Simple and intuitive API for defining HTTP routes
- Middleware support for request pre-processing and post-processing
- Chainable middleware for both global and route-specific use
- Convenient response helpers for JSON, XML, and raw data
- Support for all standard HTTP methods (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
- Bulk operations for managing multiple endpoints as a group
- Graceful shutdown support
- Minimal dependencies

## Installation

```bash
go get github.com/dbubel/intake/v2
```

## Quick Start

```go
package main

import (
    "fmt"
    "net/http"
    "time"
    
    "github.com/dbubel/intake/v2"
)

func main() {
    // Create a new Intake instance
    app := intake.New()
    
    // Define a simple handler
    helloHandler := func(w http.ResponseWriter, r *http.Request) {
        intake.Respond(w, r, http.StatusOK, []byte("Hello, World!"))
    }
    
    // Add a route
    app.AddEndpoint(http.MethodGet, "/hello", helloHandler)
    
    // Start the server with graceful shutdown
    app.Run(&http.Server{
        Addr:           ":8080",
        Handler:        app.Mux,
        ReadTimeout:    time.Second * 60,
        WriteTimeout:   time.Second * 60,
        MaxHeaderBytes: 1 << 20,
    })
}
```

## Middleware

Intake provides a flexible middleware system that allows you to intercept and process requests before they reach your handlers.

### Creating Middleware

Middleware in Intake is defined as a function that takes an `http.HandlerFunc` and returns an `http.HandlerFunc`:

```go
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Do something before the handler
        fmt.Println("Request received:", r.Method, r.URL.Path)
        
        // Call the next handler
        next(w, r)
        
        // Do something after the handler
        fmt.Println("Response sent")
    }
}
```

### Global Middleware

Global middleware is applied to all routes:

```go
app.AddGlobalMiddleware(loggingMiddleware)
```

### Route-Specific Middleware

You can add middleware to specific routes:

```go
app.AddEndpoint(http.MethodGet, "/protected", protectedHandler, authMiddleware)
```

### Middleware Chaining

Middleware can be chained together:

```go
app.AddEndpoint(http.MethodGet, "/api/data", dataHandler, 
    loggingMiddleware, 
    authMiddleware, 
    rateLimitMiddleware)
```

## Working with Endpoints

### Creating Individual Endpoints

Intake provides helper functions for creating endpoints with specific HTTP methods:

```go
// These all create endpoint objects
getEndpoint := intake.GET("/users", listUsersHandler)
postEndpoint := intake.POST("/users", createUserHandler)
putEndpoint := intake.PUT("/users/:id", updateUserHandler)
deleteEndpoint := intake.DELETE("/users/:id", deleteUserHandler)
```

### Managing Groups of Endpoints

You can group endpoints together to apply middleware to all of them:

```go
// Create a group of endpoints
apiEndpoints := intake.Endpoints{
    intake.GET("/api/users", listUsersHandler),
    intake.POST("/api/users", createUserHandler),
    intake.GET("/api/products", listProductsHandler),
}

// Apply middleware to all endpoints in the group
apiEndpoints.Use(authMiddleware)

// Add all endpoints to the app
app.AddEndpoints(apiEndpoints)
```

### Middleware Positioning

You can control the order of middleware execution:

```go
// Add middleware to the end of the chain
endpoints.Append(loggingMiddleware)

// Add middleware to the beginning of the chain (after global middleware)
endpoints.Prepend(metricMiddleware)
```

## Response Helpers

Intake provides helper functions for common response types:

```go
// JSON response
intake.RespondJSON(w, r, http.StatusOK, data)

// XML response
intake.RespondXML(w, r, http.StatusOK, data)

// Raw response
intake.Respond(w, r, http.StatusOK, []byte("Hello, World!"))
```

## OPTIONS Requests

Intake allows you to create middleware for handling OPTIONS requests for CORS:

```go
// Create a middleware for handling OPTIONS requests
optionsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Handle OPTIONS requests
        if r.Method == http.MethodOptions {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
            w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization")
            w.WriteHeader(http.StatusOK)
            return
        }
        next(w, r)
    }
}

// Add it as global middleware
app.AddGlobalMiddleware(optionsMiddleware)
```

## Complete Example

```go
package main

import (
    "fmt"
    "net/http"
    "time"
    
    "github.com/dbubel/intake/v2"
)

func main() {
    app := intake.New()
    
    // Define middleware
    loggingMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            fmt.Println("Request received:", r.Method, r.URL.Path)
            next(w, r)
            fmt.Println("Response sent")
        }
    }
    
    // Add global middleware
    app.AddGlobalMiddleware(loggingMiddleware)
    
    // Define handlers
    helloHandler := func(w http.ResponseWriter, r *http.Request) {
        intake.RespondJSON(w, r, http.StatusOK, map[string]string{
            "message": "Hello, World!",
        })
    }
    
    userHandler := func(w http.ResponseWriter, r *http.Request) {
        intake.RespondJSON(w, r, http.StatusOK, map[string]string{
            "user": "John Doe",
        })
    }
    
    // Define OPTIONS middleware
    optionsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            if r.Method == http.MethodOptions {
                w.Header().Set("Access-Control-Allow-Origin", "*")
                w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
                w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization")
                w.WriteHeader(http.StatusOK)
                return
            }
            next(w, r)
        }
    }
    
    // Add OPTIONS middleware globally
    app.AddGlobalMiddleware(optionsMiddleware)
    
    // Create endpoints
    apiEndpoints := intake.Endpoints{
        intake.GET("/api/hello", helloHandler),
        intake.GET("/api/user", userHandler),
    }
    
    // Add endpoints to the app
    app.AddEndpoints(apiEndpoints)
    
    // Add a single endpoint
    app.AddEndpoint(http.MethodGet, "/health", func(w http.ResponseWriter, r *http.Request) {
        intake.RespondJSON(w, r, http.StatusOK, map[string]string{
            "status": "healthy",
        })
    })
    
    // Print registered routes
    routes := app.GetRoutes()
    fmt.Println("Registered routes:")
    for path, methods := range routes {
        fmt.Printf("%s: %v\n", path, methods)
    }
    
    // Start the server
    fmt.Println("Server starting on :8080")
    app.Run(&http.Server{
        Addr:           ":8080",
        Handler:        app.Mux,
        ReadTimeout:    time.Second * 60,
        WriteTimeout:   time.Second * 60,
        MaxHeaderBytes: 1 << 20,
    })
}
```
## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
