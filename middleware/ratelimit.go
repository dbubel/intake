package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dbubel/intake"
	"github.com/julienschmidt/httprouter"
)

// RateLimit with limit requests that use this middleware to n requests per second
func RateLimit(n float64) func(handler intake.Handler) intake.Handler {
	var lastRequestTime = time.Now()
	return func(next intake.Handler) intake.Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			requestsPerSecond := 1 / time.Now().Sub(lastRequestTime).Seconds()
			lastRequestTime = time.Now()
			if requestsPerSecond > n {
				intake.RespondError(w, r, fmt.Errorf("too many requests"), http.StatusTooManyRequests, "rate limited")
				return
			}
			next(w, r, params)
		}
	}
}
