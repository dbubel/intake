package intake

type Endpoints []endpoint

// Use wraps a group of endpoints in middleware
func (e Endpoints) Use(mid ...MiddleWare) {
	for i := 0; i < len(e); i++ {
		e[i].MiddlewareHandlers = append(e[i].MiddlewareHandlers, mid...)
	}
}

// Append is analogous to Use but more descriptive of where the middleware is
// applied on the chain
func (e Endpoints) Append(mid ...MiddleWare) {
	e.Use(mid...)
}

// Prepend add middleware to the front of the chain for the endpoint group
func (e Endpoints) Prepend(mid ...MiddleWare) {
	mw := make([]MiddleWare, 0, 0)
	mw = append(mw, mid...)
	for i := 0; i < len(e); i++ {
		e[i].MiddlewareHandlers = append(mw, e[i].MiddlewareHandlers...)
	}
}
