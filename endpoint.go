package intake

type endpoint struct {
	Verb               string
	Path               string
	EndpointHandler    Handler
	MiddlewareHandlers []MiddleWare
}

func NewEndpoint(method, path string, endpointHandler Handler, mid ...MiddleWare) endpoint {
	return endpoint{
		Verb:               method,
		Path:               path,
		EndpointHandler:    endpointHandler,
		MiddlewareHandlers: mid,
	}
}
