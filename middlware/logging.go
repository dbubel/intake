package middlware

import (
	"github.com/bf-dbubel/intake"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type LogLevel struct {
	Log100s bool
	Log200s bool
	Log300s bool
	Log400s bool
	Log500s bool
}

func Logging(l *logrus.Logger, levels LogLevel) func(handler intake.Handler) intake.Handler {
	return func(next intake.Handler) intake.Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			t := time.Now()
			defer func() {
				var code int
				if err := intake.FromContext(r, "response-code", &code); err != nil {
					l.WithError(err).Error("error getting response code from context")
				}

				var responseLength int
				if err := intake.FromContext(r, "response-length", &responseLength); err != nil {
					l.WithError(err).Error("error getting response length from context")
				}
				printLog := func() {
					l.WithFields(logrus.Fields{
						"method":         r.Method,
						"requestUri":     r.RequestURI,
						"contentLen":     r.ContentLength,
						"responseLenBytes": responseLength,
						"responseTimeMs": time.Now().Sub(t).Milliseconds(),
						"code":           code,
					}).Info("handled request")
				}
				if code < 100 {
					printLog()
				} else if code >= 100 && code < 200 && levels.Log100s {
					printLog()
				} else if code >= 200 && code < 300 && levels.Log200s {
					printLog()
				} else if code >= 300 && code < 400 && levels.Log300s {
					printLog()
				} else if code >= 400 && code < 500 && levels.Log400s {
					printLog()
				} else if code >= 500 && levels.Log500s {
					printLog()
				}

			}()

			next(w, r, params)
		}
	}
}
