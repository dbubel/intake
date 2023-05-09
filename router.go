package intake

import (
	"net/http"
)

type Router struct {
	RadixTree *Node
}

func NewRouter() *Router {
	return &Router{
		RadixTree: NewNode(),
	}
}

func (r *Router) AddRoute(route string, method string, handler http.HandlerFunc) {
	r.RadixTree.Insert(route, method, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, found := r.RadixTree.Find(req.URL.Path, req.Method)
	if found {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}
