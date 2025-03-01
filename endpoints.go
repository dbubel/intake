// Package intake provides HTTP routing utilities.
// This file contains the Endpoints type which represents a collection of endpoints
// that can be manipulated as a group.
package intake

// Endpoints is a collection of endpoint definitions that can be manipulated as a group.
// It allows for bulk operations like adding middleware to multiple endpoints simultaneously.
// This type is useful for organizing related endpoints and applying common middleware
// or configuration to them as a unit.
type Endpoints []endpoint

// Use adds middleware to all endpoints in the collection. The middleware will be
// executed in the order they are provided, after any existing middleware.
//
// Parameters:
//   - mid: A variadic list of middleware functions to add to all endpoints
func (e Endpoints) Use(mid ...MiddleWare) {
	for i := 0; i < len(e); i++ {
		e[i].MiddlewareHandlers = append(e[i].MiddlewareHandlers, mid...)
	}
}

// Append adds middleware to the end of each endpoint's middleware chain.
// This is an alias for Use() that provides more semantic clarity about
// middleware positioning. The middleware will be executed after any existing
// middleware in each endpoint.
//
// Parameters:
//   - mid: A variadic list of middleware functions to append to all endpoints
func (e Endpoints) Append(mid ...MiddleWare) {
	e.Use(mid...)
}

// Prepend adds middleware to the beginning of each endpoint's middleware chain,
// but after any global middleware. This allows for group-specific middleware
// to execute before endpoint-specific middleware. The middleware will be executed
// in the order they are provided, before any existing endpoint-specific middleware.
//
// Parameters:
//   - mid: A variadic list of middleware functions to prepend to all endpoints
func (e Endpoints) Prepend(mid ...MiddleWare) {
	mw := make([]MiddleWare, 0)
	mw = append(mw, mid...)
	for i := 0; i < len(e); i++ {
		e[i].MiddlewareHandlers = append(mw, e[i].MiddlewareHandlers...)
	}
}
