package main

import (
	"fmt"
	"net/http"
	"strings"
)

// HandlerFunc defines the type for request handlers
type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

// node represents a node in the radix tree
type node struct {
	children map[string]*node
	param    string
	handler  HandlerFunc
}

// Router represents the HTTP router
type Router struct {
	root *node
}

// NewRouter creates a new instance of the Router
func NewRouter() *Router {
	return &Router{root: &node{children: make(map[string]*node)}}
}

// HandleFunc adds a new route and its handler to the router
func (r *Router) HandleFunc(path string, handler HandlerFunc) {
	segments := strings.Split(path, "/")
	current := r.root

	for _, segment := range segments {
		if segment == "" {
			continue
		}

		if strings.HasPrefix(segment, ":") {
			param := segment[1:]
			child, exists := current.children["param"]
			if !exists {
				child = &node{children: make(map[string]*node), param: param}
				current.children["param"] = child
			}

			current = child
		} else {
			child, exists := current.children[segment]
			if !exists {
				child = &node{children: make(map[string]*node)}
				current.children[segment] = child
			}

			current = child
		}
	}

	current.handler = handler
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	segments := strings.Split(req.URL.Path, "/")
	current := r.root
	params := make(map[string]string)

	for _, segment := range segments {
		if segment == "" {
			continue
		}

		child, exists := current.children[segment]
		if !exists {
			child, exists = current.children["param"]
			if !exists {
				http.NotFound(w, req)
				return
			}

			params[child.param] = segment
			current = child
		} else {
			current = child
		}
	}

	if current.handler != nil {
		current.handler(w, req, params)
	} else {
		http.NotFound(w, req)
	}
}

// IndexHandler handles requests to the root path
func IndexHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	fmt.Fprint(w, "Welcome to the home page!\n")
}

// HelloHandler handles requests to the /hello/:name path
func HelloHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	name := params["name"]
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

/*
func main() {
	router := NewRouter()

	// Define your routes
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/hello/:name", HelloHandler)

	// Start the HTTP server
	port := 8080
	fmt.Printf("Server listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
*/
