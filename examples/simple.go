package main

import (
	"fmt"
	middleware "github.com/dbubel/intake/middlware"
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
	intake.RespondJSON(w, r, http.StatusOK, map[string]string{
		"status": "OK",
	})
	return
}

func main() {

	apiLogger := logrus.New()
	apiLogger.SetReportCaller(true)
	apiLogger.SetLevel(logrus.DebugLevel)
	apiLogger.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	})

	app := intake.New(apiLogger)
	app.GlobalMiddleware(middleware.Logging(apiLogger, middleware.LogLevel{
		Log100s: true,
		Log200s: true,
		Log300s: true,
		Log400s: true,
		Log500s: true,
	}), middleware.Recover)
	app.AddEndpoint(http.MethodGet, "/test-get", testSimple)
	app.Run(&http.Server{
		Addr:           ":8000",
		Handler:        app.Router,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	})
}
