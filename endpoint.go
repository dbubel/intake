package intake

import "net/http"

type Endpoint struct {
	Verb               string
	Path               string
	EndpointHandler    Handler
	MiddlewareHandlers []MiddleWare
}

func NewEndpoint(method, path string, endpointHandler Handler, mid ...MiddleWare) Endpoint {
	return Endpoint{
		Verb:               method,
		Path:               path,
		EndpointHandler:    endpointHandler,
		MiddlewareHandlers: mid,
	}
}

func GET(path string, endpointHandler Handler, mid ...MiddleWare) Endpoint {
	return NewEndpoint(http.MethodGet, path, endpointHandler, mid...)
}

func POST(path string, endpointHandler Handler, mid ...MiddleWare) Endpoint {
	return NewEndpoint(http.MethodPost, path, endpointHandler, mid...)
}

func PUT(path string, endpointHandler Handler, mid ...MiddleWare) Endpoint {
	return NewEndpoint(http.MethodPut, path, endpointHandler, mid...)
}

func DELETE(path string, endpointHandler Handler, mid ...MiddleWare) Endpoint {
	return NewEndpoint(http.MethodDelete, path, endpointHandler, mid...)
}

func PATCH(path string, endpointHandler Handler, mid ...MiddleWare) Endpoint {
	return NewEndpoint(http.MethodPatch, path, endpointHandler, mid...)
}

func HEAD(path string, endpointHandler Handler, mid ...MiddleWare) Endpoint {
	return NewEndpoint(http.MethodHead, path, endpointHandler, mid...)
}
