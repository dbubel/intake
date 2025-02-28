// Package intake provides HTTP routing utilities.
package intake

import "net/http"

// endpoint represents a single HTTP route with its associated handler and middleware.
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
func NewEndpoint(method, path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return endpoint{
		Verb:               method,
		Path:               path,
		EndpointHandler:    endpointHandler,
		MiddlewareHandlers: mid,
	}
}

// GET creates a new endpoint for handling HTTP GET requests at the specified path.
func GET(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return NewEndpoint(http.MethodGet, path, endpointHandler, mid...)
}

// POST creates a new endpoint for handling HTTP POST requests at the specified path.
func POST(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return NewEndpoint(http.MethodPost, path, endpointHandler, mid...)
}

// PUT creates a new endpoint for handling HTTP PUT requests at the specified path.
func PUT(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return NewEndpoint(http.MethodPut, path, endpointHandler, mid...)
}

// DELETE creates a new endpoint for handling HTTP DELETE requests at the specified path.
func DELETE(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return NewEndpoint(http.MethodDelete, path, endpointHandler, mid...)
}

// PATCH creates a new endpoint for handling HTTP PATCH requests at the specified path.
func PATCH(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return NewEndpoint(http.MethodPatch, path, endpointHandler, mid...)
}

// HEAD creates a new endpoint for handling HTTP HEAD requests at the specified path.
func HEAD(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) endpoint {
	return NewEndpoint(http.MethodHead, path, endpointHandler, mid...)
}
