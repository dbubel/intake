package intake

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Intake struct {
	Router           *Router
	GlobalMiddleware []MiddleWare
	OptionsHandler   http.HandlerFunc
	NotFoundHandler  http.HandlerFunc // TODO:(dean) implement
}

// New Create a new intake application
func New() *Intake {
	return &Intake{
		Router:           NewRouter(),
		GlobalMiddleware: make([]MiddleWare, 0, 0),
	}
}

// AddGlobalMiddleware Global middleware MUST be added before other routes.
// adding a global middleware will cause that middleware to execute for all
// routes that are part of the router.
func (a *Intake) AddGlobalMiddleware(mw MiddleWare) {
	a.GlobalMiddleware = append(a.GlobalMiddleware, mw)
}

// AddEndpoints adds a list of Endpoints to the router
func (a *Intake) AddEndpoints(e ...Endpoints) {
	for x := 0; x < len(e); x++ {
		for i := 0; i < len(e[x]); i++ {
			a.AddEndpoint(e[x][i].Path, e[x][i].Verb, e[x][i].EndpointHandler, e[x][i].MiddlewareHandlers...)
		}
	}
}

// AddEndpoint adds an endpoint handler to the router.
func (a *Intake) AddEndpoint(path string, verb string, finalHandler http.HandlerFunc, middleware ...MiddleWare) {
	// Prepend the global middlewares to the route specific middleware
	// global middleware will be called first in the chain in the order they are added
	mws := append(a.GlobalMiddleware, middleware...)
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] != nil {
			finalHandler = mws[i](finalHandler)
		}
	}

	// add our wrapped handler chain to the router
	a.Router.AddRoute(path, verb, finalHandler)

	// if we want to handle OPTIONS request add the handler to router
	if a.OptionsHandler != nil {
		a.Router.AddRoute(path, http.MethodOptions, a.OptionsHandler)
	}

	log.Println("added route", verb, path)
}

// Run starts the http server. Waits for a signal to gracefully shutdown.
func (a *Intake) Run(server *http.Server) {
	serverErrors := make(chan error, 1)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	log.Println("server started on port", server.Addr)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Println("error starting server", err.Error())
	case <-osSignals:
		log.Println("shutdown received shedding connections...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Println("graceful shutdown did not complete in allowed time", err.Error())
			if err := server.Close(); err != nil {
				log.Println("could not stop http server", err.Error())
			}
		}
		log.Println("shutdown OK")
	}
}
