package middlware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dbubel/intake"
	"github.com/julienschmidt/httprouter"
)

// RateLimitIP with limit requests that use this middleware to n requests per second
func RateLimitIP(n float64) func(handler intake.Handler) intake.Handler {
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
	return func(next intake.Handler) intake.Handler {
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
				intake.RespondError(w, r, fmt.Errorf("too many requests"), http.StatusTooManyRequests, "rate limited")
				return
			}
			next(w, r, params)
		}
	}
}
