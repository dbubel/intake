package intake

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
	"time"
)

func (a *Intake) Logging(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		t := time.Now()
		defer func() {
			var code int
			if err := FromContext(r, "response-code", &code); err != nil {
				a.logger.WithError(err).Error("error getting response code from context")
			}

			var responseLength int
			if err := FromContext(r, "response-length", &responseLength); err != nil {
				a.logger.WithError(err).Error("error getting response length from context")
			}

			a.logger.WithFields(logrus.Fields{
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

// Recover needs to be the first middleware in chain
func (a *Intake) Recover(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				RespondError(w, r, fmt.Errorf("error panic"), http.StatusInternalServerError, "server recovered from a panic")
			}
		}()
		next(w, r, params)
	}
}

// RateLimit will limit requests that use this middleware to n requests per second
func (a *Intake) RateLimit(n float64) func(handler Handler) Handler {
	var lastRequestTime = time.Now()
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			requestsPerSecond := 1 / time.Now().Sub(lastRequestTime).Seconds()
			lastRequestTime = time.Now()
			if requestsPerSecond > n {
				RespondError(w, r, fmt.Errorf("too many requests"), http.StatusTooManyRequests, "rate limited")
				return
			}
			next(w, r, params)
		}
	}
}

// RateLimitIP limit the number of requests per second per IP
func (a *Intake) RateLimitIP(n float64) func(handler Handler) Handler {
	var ipMap map[string]time.Time
	ipMap = make(map[string]time.Time)

	fn := func(s string) string {
		ipSplit := strings.Split(s, ":")
		if len(ipSplit) > 0 {
			ip := net.ParseIP(ipSplit[0])
			if ip != nil {
				return ip.String()
			}
		}

		return "noip"
	}
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			var requestsPerSecond float64
			requestIp := fn(r.RemoteAddr)
			val, exists := ipMap[requestIp]
			if exists != true {
				requestsPerSecond = 0
			} else {
				requestsPerSecond = 1 / time.Now().Sub(val).Seconds()
			}

			ipMap[requestIp] = time.Now()
			if requestsPerSecond > n {
				RespondError(w, r, fmt.Errorf("too many requests"), http.StatusTooManyRequests, "rate limited")
				return
			}
			next(w, r, params)
		}
	}
}

// Timeout if added will created a context that will cancel after t. This cancel
// will affect all downstream uses of the context
func (a *Intake) Timeout(t time.Duration) func(handler Handler) Handler {
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			// Create a context that is both manually cancellable and will signal
			// cancel at the specified duration.
			ctx, cancel := context.WithTimeout(r.Context(), t)
			defer cancel()
			*r = *r.WithContext(ctx)
			next(w, r, params)
		}
	}
}
