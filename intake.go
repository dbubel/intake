package intake

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type Handler func(w http.ResponseWriter, r *http.Request, params httprouter.Params)
type MiddleWare func(Handler) Handler

type Intake struct {
	Router   *httprouter.Router
	logger   *logrus.Logger
	GlobalMw []MiddleWare
}

func New(log *logrus.Logger) *Intake {
	return &Intake{
		Router:   httprouter.New(),
		logger:   log,
		GlobalMw: make([]MiddleWare, 0, 0),
	}
}

func NewDefault() *Intake {
	apiLogger := logrus.New()
	apiLogger.SetReportCaller(true)
	apiLogger.SetLevel(logrus.DebugLevel)
	apiLogger.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			//filename := path.Base(f.File)
			return fmt.Sprintf("%s", f.Function), fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	})
	return &Intake{
		Router: httprouter.New(),
		logger: apiLogger,
	}
}

func (a *Intake) AddGlobal(mw MiddleWare) {
	a.GlobalMw = append(a.GlobalMw, mw)
}

func (a *Intake) AddEndpoints(e ...Endpoints) {
	for x := 0; x < len(e); x++ {
		for i := 0; i < len(e[x]); i++ {
			a.AddEndpoint(e[x][i].Verb, e[x][i].Path, e[x][i].EndpointHandler, e[x][i].MiddlewareHandlers...)
		}
	}
}

func (a *Intake) AddEndpoint(verb string, path string, finalHandler Handler, middleware ...MiddleWare) {
	// Prepend the global middlewares to the route specific middleware
	// global middleware will be called first in the chain in the order they are added
	mws := append(a.GlobalMw, middleware...)
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] != nil {
			finalHandler = mws[i](finalHandler)
		}
	}

	// Our wrapped function chain in a compatible httprouter AddEndpoint func
	a.Router.Handle(verb, path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		finalHandler(w, r, params)
	})
	a.logger.WithFields(logrus.Fields{"verb": verb, "path": path}).Debug("added route")
}

func (a *Intake) Run(server *http.Server) {
	serverErrors := make(chan error, 1)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	a.logger.WithFields(logrus.Fields{"addr": server.Addr}).Info("server starting")

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		a.logger.WithError(err).Error("error starting server")
	case <-osSignals:
		a.logger.Info("shutdown received shedding connections...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			a.logger.WithError(err).Error("graceful shutdown did not complete in allowed time")
			if err := server.Close(); err != nil {
				a.logger.WithError(err).Error("could not stop http server")
			}
		}
		a.logger.Info("shutdown OK")
	}
}
