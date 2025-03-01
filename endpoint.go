// Package intake provides HTTP routing utilities.
// This file contains the endpoint definition and convenience functions for
// creating endpoints for different HTTP methods.
package intake

import "net/http"

// endpoint represents a single HTTP route with its associated handler and middleware.
// An endpoint encapsulates all the information needed to register and handle
// a specific HTTP route, including its HTTP method, path, handler function,
// and any middleware specific to this endpoint.
type endpoint struct {
	// Verb is the HTTP method (GET, POST, etc.)
	Verb string
	// Path is the URL path for this endpoint
	Path string
	// EndpointHandler is the main handler function for this endpoint
	EndpointHandler http.HandlerFunc
	// MiddlewareHandlers are the middleware functions specific to this endpoint
	MiddlewareHandlers []MiddleWare
}

// NewEndpoint creates a new endpoint with the specified HTTP method, path, handler, and optional middleware.
// This is the general constructor function for creating endpoints. For convenience,
// method-specific constructors (GET, POST, etc.) are also provided.
//
// Parameters:
//   - method: The HTTP method (GET, POST, PUT, etc.)
//   - path: The URL path for this endpoint
//   - endpointHandler: The handler function for this endpoint
//   - mid: Optional middleware functions specific to this endpoint
//
// Returns:
//   - A new endpoint instance configured with the provided parameters
func NewEndpoint(method, path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return endpoint{
		Verb:               method,
		Path:               path,
		EndpointHandler:    endpointHandler,
		MiddlewareHandlers: mid,
	}
}

// GET creates a new endpoint for handling HTTP GET requests at the specified path.
// This is a convenience function that calls NewEndpoint with http.MethodGet as the method.
//
// Parameters:
//   - path: The URL path for this endpoint
//   - endpointHandler: The handler function for this endpoint
//   - mid: Optional middleware functions specific to this endpoint
//
// Returns:
//   - A new endpoint instance configured for GET requests
func GET(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return NewEndpoint(http.MethodGet, path, endpointHandler, mid...)
}

// POST creates a new endpoint for handling HTTP POST requests at the specified path.
// This is a convenience function that calls NewEndpoint with http.MethodPost as the method.
//
// Parameters:
//   - path: The URL path for this endpoint
//   - endpointHandler: The handler function for this endpoint
//   - mid: Optional middleware functions specific to this endpoint
//
// Returns:
//   - A new endpoint instance configured for POST requests
func POST(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return NewEndpoint(http.MethodPost, path, endpointHandler, mid...)
}

// PUT creates a new endpoint for handling HTTP PUT requests at the specified path.
// This is a convenience function that calls NewEndpoint with http.MethodPut as the method.
//
// Parameters:
//   - path: The URL path for this endpoint
//   - endpointHandler: The handler function for this endpoint
//   - mid: Optional middleware functions specific to this endpoint
//
// Returns:
//   - A new endpoint instance configured for PUT requests
func PUT(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return NewEndpoint(http.MethodPut, path, endpointHandler, mid...)
}

// DELETE creates a new endpoint for handling HTTP DELETE requests at the specified path.
// This is a convenience function that calls NewEndpoint with http.MethodDelete as the method.
//
// Parameters:
//   - path: The URL path for this endpoint
//   - endpointHandler: The handler function for this endpoint
//   - mid: Optional middleware functions specific to this endpoint
//
// Returns:
//   - A new endpoint instance configured for DELETE requests
func DELETE(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return NewEndpoint(http.MethodDelete, path, endpointHandler, mid...)
}

// PATCH creates a new endpoint for handling HTTP PATCH requests at the specified path.
// This is a convenience function that calls NewEndpoint with http.MethodPatch as the method.
//
// Parameters:
//   - path: The URL path for this endpoint
//   - endpointHandler: The handler function for this endpoint
//   - mid: Optional middleware functions specific to this endpoint
//
// Returns:
//   - A new endpoint instance configured for PATCH requests
func PATCH(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return NewEndpoint(http.MethodPatch, path, endpointHandler, mid...)
}

// HEAD creates a new endpoint for handling HTTP HEAD requests at the specified path.
// This is a convenience function that calls NewEndpoint with http.MethodHead as the method.
//
// Parameters:
//   - path: The URL path for this endpoint
//   - endpointHandler: The handler function for this endpoint
//   - mid: Optional middleware functions specific to this endpoint
//
// Returns:
//   - A new endpoint instance configured for HEAD requests
func HEAD(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return NewEndpoint(http.MethodHead, path, endpointHandler, mid...)
}
