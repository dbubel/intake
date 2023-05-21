package intake

import (
	"fmt"
	"net/http"
)

type Router struct {
	routes map[string]http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]http.HandlerFunc),
	}
}

func (r *Router) AddRoute(route string, method string, handler http.HandlerFunc) {
	fmt.Println(method+route)
	r.routes[method+route] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method+req.URL.Path)
	handler, found := r.routes[req.Method+req.URL.Path]
	if found {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}
