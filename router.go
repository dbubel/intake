package intake

import (
	"net/http"
	"strings"
)

type Router struct {
	routes                map[string]http.HandlerFunc
	RedirectTrailingSlash bool
}

func NewRouter() *Router {
	return &Router{
		routes:                make(map[string]http.HandlerFunc),
		RedirectTrailingSlash: true,
	}
}

func (r *Router) AddRoute(route string, method string, handler http.HandlerFunc) {
	r.routes[method+strings.TrimSuffix(route, "/")] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Redirect if path ends with a trailing slash
	if r.RedirectTrailingSlash {
		if RedirectTrailingSlash(w, req) {
			return
		}
	}

	handler, found := r.routes[req.Method+req.URL.Path]
	if found {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}

// RedirectTrailingSlash redirects paths with a trailing slash to the same path, minus the slash.
func RedirectTrailingSlash(w http.ResponseWriter, req *http.Request) bool {
	path := req.URL.Path
	if path != "/" && strings.HasSuffix(path, "/") {
		req.URL.Path = strings.TrimSuffix(path, "/")
		http.Redirect(w, req, req.URL.String(), http.StatusMovedPermanently)
		return true
	}
	return false
}
