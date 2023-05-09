package intake

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Intake struct {
	Router           *Router
	GlobalMiddleware []MiddleWare
}

func NewDefault() *Intake {
	return &Intake{
		Router: NewRouter(),
		GlobalMiddleware: make([]MiddleWare, 0, 0),
	}
}

// AddGlobalMiddleware Global middleware MUST be added before other routes
func (a *Intake) AddGlobalMiddleware(mw MiddleWare) {
	a.GlobalMiddleware = append(a.GlobalMiddleware, mw)
}

func (a *Intake) AddEndpoints(e ...Endpoints) {
	for x := 0; x < len(e); x++ {
		for i := 0; i < len(e[x]); i++ {
			a.AddEndpoint( e[x][i].Path, e[x][i].Verb,e[x][i].EndpointHandler, e[x][i].MiddlewareHandlers...)
		}
	}
}


func (a *Intake) AddEndpoint(path string, verb string, finalHandler http.HandlerFunc, middleware ...MiddleWare) {
	// Prepend the global middlewares to the route specific middleware
	// global middleware will be called first in the chain in the order they are added
	mws := append(a.GlobalMiddleware, middleware...)
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] != nil {
			finalHandler = mws[i](finalHandler)
		}
	}

	// Our wrapped function chain in a compatible httprouter AddEndpoint func
	a.Router.AddRoute(path, verb,finalHandler)
	//a.Logger.WithFields(logrus.Fields{"verb": verb, "path": path}).Debug("added route")
}

func (a *Intake) Run(server *http.Server) {
	serverErrors := make(chan error, 1)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	//a.Logger.WithFields(logrus.Fields{"addr": server.Addr}).Info("server starting")

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		_=err
		//a.Logger.WithError(err).Error("error starting server")
	case <-osSignals:
		//a.Logger.Info("shutdown received shedding connections...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			//a.Logger.WithError(err).Error("graceful shutdown did not complete in allowed time")
			if err := server.Close(); err != nil {
				//a.Logger.WithError(err).Error("could not stop http server")
			}
		}
		//a.Logger.Info("shutdown OK")
	}
}
