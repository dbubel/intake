package intake

import (
	"net/http"
	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
)

// AttachPprofTraceEndpoints attaches with index page for viewing
func (a *Intake) AttachPprofTraceEndpoints() {
	a.Logger.Info("attaching debug pprof endpoints")
	a.Router.Handler(http.MethodGet, "/debug/pprof/*item", http.DefaultServeMux)
}

// DebugTraceEndpoints follows the same middleware route pattern without the pprof
// index page
func DebugTraceEndpoints(mw ...MiddleWare) Endpoints {
	endpoints := Endpoints{
		GET("/debug/pprof/cmdline", DEBUGPProfCmdLine),
		GET("/debug/pprof/profile", DEBUGPProfProfile),
		GET("/debug/pprof/symbol", DEBUGPProfSymbol),
		GET("/debug/pprof/trace", DEBUGPProfTrace),
	}
	endpoints.Use(mw...)
	return endpoints
}

func DEBUGPProfCmdLine(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Cmdline(w, r)
}

func DEBUGPProfProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Profile(w, r)
}

func DEBUGPProfSymbol(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Symbol(w, r)
}

func DEBUGPProfTrace(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pprof.Trace(w, r)
}
