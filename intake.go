// Package intake implements a simple HTTP router with middleware support.
// It provides a lightweight framework for building HTTP services with support for
// middleware, panic recovery, and graceful shutdown.
package intake

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// MiddleWare defines a function that wraps an http.HandlerFunc with additional behavior.
// Middleware functions can be used to add cross-cutting concerns such as logging,
// authentication, request tracing, or any other functionality that should be applied
// to multiple endpoints.
type MiddleWare func(http.HandlerFunc) http.HandlerFunc

// Intake represents an HTTP router with middleware and panic recovery support.
// It provides a simple API for registering routes, applying middleware, and
// handling HTTP requests. The Intake struct encapsulates all the functionality
// needed to build and run an HTTP service.
type Intake struct {
	// Mux is the underlying HTTP request multiplexer
	Mux *http.ServeMux
	// PanicHandler handles any panics that occur during request processing
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
	// GlobalMiddleware contains middleware applied to all routes
	GlobalMiddleware []MiddleWare
	// OptionsHandlerFunc handles OPTIONS requests
	OptionsHandlerFunc http.HandlerFunc
	// optionsPaths tracks paths with OPTIONS handlers
	optionsPaths map[string]bool
	// registeredRoutes maps paths to their HTTP methods
	registeredRoutes map[string][]string
}

// New creates a new Intake instance with initialized maps and slices.
// It sets up the internal data structures needed for routing and middleware
// management. This function should be called to create a new router before
// registering any routes or middleware.
func New() *Intake {
	return &Intake{
		GlobalMiddleware: make([]MiddleWare, 0),
		Mux:              http.NewServeMux(),
		optionsPaths:     make(map[string]bool),
		registeredRoutes: make(map[string][]string),
	}
}

// AddGlobalMiddleware adds middleware that will be applied to all routes.
// Global middleware must be added before registering routes. The middleware
// functions are executed in the order they are added, with the first added
// middleware being the outermost wrapper around the handler function.
//
// Parameters:
//   - mw: The middleware function to add to the global middleware chain.
func (a *Intake) AddGlobalMiddleware(mw MiddleWare) {
	a.GlobalMiddleware = append(a.GlobalMiddleware, mw)
}

// AddEndpoints registers multiple endpoints at once.
// This is a convenience method that allows registering multiple endpoints
// in a single call, which can improve code readability when setting up
// multiple related routes.
//
// Parameters:
//   - e: A variadic parameter of Endpoints slices to register.
func (a *Intake) AddEndpoints(e ...Endpoints) {
	for x := 0; x < len(e); x++ {
		for i := 0; i < len(e[x]); i++ {
			a.AddEndpoint(e[x][i].Verb, e[x][i].Path, e[x][i].EndpointHandler, e[x][i].MiddlewareHandlers...)
		}
	}
}

// OptionsHandler sets the handler for OPTIONS requests across all routes.
// This handler will be automatically applied to all routes that have been
// registered. It's particularly useful for implementing CORS support by
// responding to preflight requests.
//
// Parameters:
//   - h: The handler function to use for OPTIONS requests.
func (a *Intake) OptionsHandler(h http.HandlerFunc) {
	a.OptionsHandlerFunc = h
	a.optionsPaths = make(map[string]bool)
}

// AddEndpoint registers a new route with the specified HTTP method and path.
// It applies both global and route-specific middleware to the handler. The
// middleware is applied in order, with global middleware being applied before
// route-specific middleware.
//
// Parameters:
//   - verb: The HTTP method (GET, POST, PUT, DELETE, etc.)
//   - path: The URL path to register the handler for
//   - finalHandler: The handler function that will process the request
//   - middleware: Optional route-specific middleware functions
func (a *Intake) AddEndpoint(verb string, path string, finalHandler http.HandlerFunc, middleware ...MiddleWare) {
	mws := append(a.GlobalMiddleware, middleware...)
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] != nil {
			finalHandler = mws[i](finalHandler)
		}
	}

	// Store the route in our registry
	if methods, exists := a.registeredRoutes[path]; exists {
		a.registeredRoutes[path] = append(methods, verb)
	} else {
		a.registeredRoutes[path] = []string{verb}
	}

	handlerKey := fmt.Sprintf("%s %s", verb, path)
	a.Mux.HandleFunc(handlerKey, func(w http.ResponseWriter, r *http.Request) {
		finalHandler(w, r)
	})

	if a.OptionsHandlerFunc != nil && !a.optionsPaths[path] {
		optionsKey := fmt.Sprintf("%s %s", http.MethodOptions, path)
		a.Mux.HandleFunc(optionsKey, a.OptionsHandlerFunc)
		if methods, exists := a.registeredRoutes[path]; exists {
			a.registeredRoutes[path] = append(methods, http.MethodOptions)
		} else {
			a.registeredRoutes[path] = []string{http.MethodOptions}
		}
		a.optionsPaths[path] = true
	}
}

// Run starts the HTTP server and handles graceful shutdown on SIGINT/SIGTERM.
// This method blocks until the server is shut down either by an error or by
// receiving a termination signal. When a signal is received, the server attempts
// to gracefully shut down, allowing in-flight requests to complete within a
// timeout period.
//
// Parameters:
//   - server: The configured http.Server instance to run
func (a *Intake) Run(server *http.Server) {
	serverErrors := make(chan error, 1)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	// Blocking main and waiting for shutdown.
	select {
	case <-serverErrors:
	case <-osSignals:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			if err := server.Close(); err != nil {
			}
		}
	}
}

// GetRoutes returns a map of paths to their supported HTTP methods.
// This can be useful for debugging or for generating documentation about
// the available endpoints. The returned map is a copy of the internal
// route registry, so modifications to it will not affect the router.
//
// Returns:
//   - A map where keys are URL paths and values are slices of HTTP methods
//     supported by each path.
func (a *Intake) GetRoutes() map[string][]string {
	routes := make(map[string][]string)
	for path, methods := range a.registeredRoutes {
		routes[path] = append([]string{}, methods...)
	}
	return routes
}
