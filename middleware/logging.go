package middleware

import (
	"github.com/dbubel/intake"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)


func Logging(l *logrus.Logger) func(handler intake.Handler) intake.Handler {
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

				l.WithFields(logrus.Fields{
					"method":           r.Method,
					"requestUri":       r.RequestURI,
					"contentLen":       r.ContentLength,
					"responseLenBytes": responseLength,
					"responseTimeMs":   time.Now().Sub(t).Milliseconds(),
					"code":             code,
				}).Info("handled request")

			}()

			next(w, r, params)
		}
	}
}
