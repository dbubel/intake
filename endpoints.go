package intake

type Endpoints []endpoint

// Use wraps a group of endpoints in middleware
func (e Endpoints) Use(mid ...MiddleWare) {
	for i := 0; i < len(e); i++ {
		e[i].MiddlewareHandlers = append(e[i].MiddlewareHandlers, mid...)
	}
}
