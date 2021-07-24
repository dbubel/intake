package main

import (
	"fmt"
	"github.com/dbubel/intake/middleware"
	"net/http"
	"runtime"
	"time"

	"github.com/dbubel/intake"
	"github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

type test struct {
	Name string `validate:"required"`
	Addr string `validate:"required"`
}

func testSimple(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("handle")
	defer func() {
		var s int
		intake.FromContext(r, "response-code", &s)
		fmt.Printf("handle leave [%d]\n", s)
	}()
	intake.RespondJSON(w, r, http.StatusOK, map[string]string{
		"status": "OK",
	})
	return
}

func mw1(next intake.Handler) intake.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		fmt.Println("m1")
		defer func() {
			var s int
			intake.FromContext(r, "response-code", &s)
			fmt.Printf("mw1 leave [%d]\n", s)
		}()
		next(w, r, params)
	}
}

func mw2(next intake.Handler) intake.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		fmt.Println("m2")
		defer func() {

			var s int
			intake.FromContext(r, "response-code", &s)
			fmt.Printf("mw2 leave [%d]\n", s)
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
			return "", fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	})

	app := intake.New(apiLogger)
	eps := intake.Endpoints{
		intake.NewEndpoint(http.MethodGet, "/test-get", testSimple),
	}

	loggingMw := middleware.Logging(apiLogger, middleware.LogLevel{
		Log100s: true,
		Log200s: true,
		Log300s: true,
		Log400s: true,
		Log500s: true,
	})


	eps.Use(loggingMw)
	app.AddEndpoints(eps)

	app.Run(&http.Server{
		Addr:           ":8000",
		Handler:        app.Router,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	})
}
