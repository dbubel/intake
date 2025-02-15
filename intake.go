package intake

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type MiddleWare func(http.HandlerFunc) http.HandlerFunc

type Intake struct {
	Mux                *http.ServeMux
	PanicHandler       func(http.ResponseWriter, *http.Request, interface{})
	GlobalMiddleware   []MiddleWare
	OptionsHandlerFunc http.HandlerFunc
	// Track paths that already have OPTIONS handlers
	optionsPaths map[string]bool
}

func New() *Intake {
	return &Intake{
		GlobalMiddleware: make([]MiddleWare, 0),
		Mux:              http.NewServeMux(),
		optionsPaths:     make(map[string]bool),
	}
}

// AddGlobalMiddleware Global middleware MUST be added before other routes
func (a *Intake) AddGlobalMiddleware(mw MiddleWare) {
	a.GlobalMiddleware = append(a.GlobalMiddleware, mw)
}

func (a *Intake) AddEndpoints(e ...Endpoints) {
	for x := 0; x < len(e); x++ {
		for i := 0; i < len(e[x]); i++ {
			a.AddEndpoint(e[x][i].Verb, e[x][i].Path, e[x][i].EndpointHandler, e[x][i].MiddlewareHandlers...)
		}
	}
}

func (a *Intake) OptionsHandler(h http.HandlerFunc) {
	a.OptionsHandlerFunc = h
	// When setting a new OPTIONS handler, we need to clear existing paths
	// in case the handler has changed
	a.optionsPaths = make(map[string]bool)
}

func (a *Intake) AddEndpoint(verb string, path string, finalHandler http.HandlerFunc, middleware ...MiddleWare) {
	// Prepend the global middlewares to the route specific middleware
	// global middleware will be called first in the chain in the order they are added
	mws := append(a.GlobalMiddleware, middleware...)
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] != nil {
			finalHandler = mws[i](finalHandler)
		}
	}

	a.Mux.HandleFunc(fmt.Sprintf("%s %s", verb, path), func(w http.ResponseWriter, r *http.Request) {
		finalHandler(w, r)
	})

	// Only add OPTIONS handler if we have one and haven't already added it for this path
	if a.OptionsHandlerFunc != nil && !a.optionsPaths[path] {
		a.Mux.HandleFunc(fmt.Sprintf("%s %s", http.MethodOptions, path), a.OptionsHandlerFunc)
		a.optionsPaths[path] = true
	}

	fmt.Printf("added route %s %s\n", verb, path)
}

func (a *Intake) Run(server *http.Server) {
	serverErrors := make(chan error, 1)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	fmt.Printf("server running on port [%s]\n", server.Addr)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		fmt.Printf("error starting server [%s]\n", err.Error())
	case <-osSignals:
		fmt.Println("shutdown received shedding connections...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			fmt.Println("graceful shutdown did not complete in allowed time")
			if err := server.Close(); err != nil {
				fmt.Printf("error calling close for server shut down [%s]\n", err.Error())
			}
		}
		fmt.Println("server shutdown OK")
	}
}

func (r *Intake) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, rcv)
	}
}
