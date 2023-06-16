package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/dbubel/intake"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type test struct {
	Name string `validate:"required"`
	Addr string `validate:"required"`
}

func testSimple(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	intake.RespondJSON(w, r, http.StatusOK, map[string]string{
		"status": "OK",
	})
	return
}

func mw1(next intake.Handler) intake.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		defer func() {
			var s int
			intake.FromContext(r, "response-code", &s)
		}()
		next(w, r, params)
	}
}

func mw2(next intake.Handler) intake.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		defer func() {
			var s int
			intake.FromContext(r, "response-code", &s)
		}()
		next(w, r, params)
	}
}

func main() {
	apiLogger := logrus.New()
	apiLogger.SetLevel(logrus.DebugLevel)
	apiLogger.SetReportCaller(true)
	apiLogger.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.File, "/")[len(strings.Split(f.File, "/"))-1]
			return "", fmt.Sprintf("%s:%d", s, f.Line)
		},
	})

	app := intake.New(apiLogger)
	eps := intake.Endpoints{
		intake.GET("/test-get", testSimple),
	}

	mw := Middleware{logger: apiLogger}

	app.AddGlobalMiddleware(mw.Logging)
	app.AddGlobalMiddleware(mw2)
	app.AddGlobalMiddleware(mw1)
	app.AddEndpoints(eps)

	app.AttachPprofTraceEndpoints()

	app.Run(&http.Server{
		Addr:           ":8000",
		Handler:        app.Router,
		ReadTimeout:    time.Second * 30,
		WriteTimeout:   time.Second * 30,
		MaxHeaderBytes: 1 << 20,
	})
}
