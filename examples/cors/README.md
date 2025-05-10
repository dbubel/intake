# CORS Support in Intake

This example demonstrates how to configure and use CORS (Cross-Origin Resource Sharing) support in the Intake HTTP router.

## What is CORS?

CORS is a security feature implemented by browsers that restricts web pages from making requests to a different domain than the one that served the original page. This is a security measure to prevent malicious websites from making unauthorized requests to APIs or services on behalf of the user.

To allow cross-origin requests, the server needs to include specific HTTP headers in its responses. Intake provides a CORS middleware that handles these headers automatically.

## Features

The CORS middleware in Intake supports:

- Configurable allowed origins (specific domains or wildcard `*` for all domains)
- Configurable allowed HTTP methods
- Configurable allowed headers
- Support for credentials (cookies, HTTP auth)
- Preflight request handling (OPTIONS requests)
- Configurable max age for preflight caching

## How to Use

### Basic Usage

```go
// Create a new Intake router
router := intake.New()

// Use default CORS settings (allow all origins)
corsMiddleware := intake.CORS(intake.DefaultCORSConfig())

// Add CORS middleware globally to all routes
router.AddGlobalMiddleware(corsMiddleware)

// Don't forget to add OPTIONS endpoints for all routes
// This should be done after registering all your routes
router.AddOptionsEndpoints()
```

### Custom CORS Configuration

You can customize the CORS settings to fit your security requirements:

```go
// Create custom CORS config
corsConfig := intake.CORSConfig{
    AllowedOrigins:   []string{"https://example.com", "https://api.example.com"},
    AllowedMethods:   []string{http.MethodGet, http.MethodPost},
    AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           3600, // 1 hour
}

// Create middleware with custom config
corsMiddleware := intake.CORS(corsConfig)

// Apply to routes
router.AddGlobalMiddleware(corsMiddleware)
// or to specific endpoint groups:
endpoints := intake.Endpoints{
    intake.GET("/api/v1/data", dataHandler),
    intake.POST("/api/v1/data", createDataHandler),
}
endpoints.Use(corsMiddleware)
```

## Testing CORS

This example includes a test HTML file that you can use to test CORS functionality with your Intake server:

1. Start the example server:
   ```
   go run main.go
   ```

2. Open the `test.html` file in a web browser or host it on a different port.

3. Click the buttons to test different types of requests and see the CORS headers in the response.

## Browser Support

The CORS implementation in Intake follows the W3C specification and should work with all modern browsers.