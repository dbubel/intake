// Package intake implements a simple HTTP router with middleware support.
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
type MiddleWare func(http.HandlerFunc) http.HandlerFunc

// Intake represents an HTTP router with middleware and panic recovery support.
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
func New() *Intake {
	return &Intake{
		GlobalMiddleware: make([]MiddleWare, 0),
		Mux:              http.NewServeMux(),
		optionsPaths:     make(map[string]bool),
		registeredRoutes: make(map[string][]string),
	}
}

// AddGlobalMiddleware adds middleware that will be applied to all routes.
// Global middleware must be added before registering routes.
func (a *Intake) AddGlobalMiddleware(mw MiddleWare) {
	a.GlobalMiddleware = append(a.GlobalMiddleware, mw)
}

// AddEndpoints registers multiple endpoints at once.
func (a *Intake) AddEndpoints(e ...Endpoints) {
	for x := 0; x < len(e); x++ {
		for i := 0; i < len(e[x]); i++ {
			a.AddEndpoint(e[x][i].Verb, e[x][i].Path, e[x][i].EndpointHandler, e[x][i].MiddlewareHandlers...)
		}
	}
}

// OptionsHandler sets the handler for OPTIONS requests across all routes.
func (a *Intake) OptionsHandler(h http.HandlerFunc) {
	a.OptionsHandlerFunc = h
	a.optionsPaths = make(map[string]bool)
}

// AddEndpoint registers a new route with the specified HTTP method and path.
// It applies both global and route-specific middleware to the handler.
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
	// comment

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
func (a *Intake) GetRoutes() map[string][]string {
	routes := make(map[string][]string)
	for path, methods := range a.registeredRoutes {
		routes[path] = append([]string{}, methods...)
	}
	return routes
}
