package intake

import (
	"fmt"
	"net/http"
)

type MiddleWare func(http.HandlerFunc) http.HandlerFunc
type Endpoint struct {
	Verb               string
	Path               string
	EndpointHandler    http.HandlerFunc
	MiddlewareHandlers []MiddleWare
}

func NewEndpoint(method, path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) *Endpoint {
	return &Endpoint{
		Verb:               method,
		Path:               path,
		EndpointHandler:    endpointHandler,
		MiddlewareHandlers: mid,
	}
}

func (e *Endpoint) Prefix(prefix string) {
	e.Path = fmt.Sprintf("%s%s", prefix, e.Path)
}

func GET(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) *Endpoint {
	return NewEndpoint(http.MethodGet, path, endpointHandler, mid...)
}

func POST(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) *Endpoint {
	return NewEndpoint(http.MethodPost, path, endpointHandler, mid...)
}

func PUT(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) *Endpoint {
	return NewEndpoint(http.MethodPut, path, endpointHandler, mid...)
}

func DELETE(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) *Endpoint {
	return NewEndpoint(http.MethodDelete, path, endpointHandler, mid...)
}

func PATCH(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) *Endpoint {
	return NewEndpoint(http.MethodPatch, path, endpointHandler, mid...)
}

func HEAD(path string, endpointHandler http.HandlerFunc, mid ...MiddleWare) *Endpoint {
	return NewEndpoint(http.MethodHead, path, endpointHandler, mid...)
}
