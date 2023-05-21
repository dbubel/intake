package main

// These are example middlewares
// These are not production tested. Please do your own testing of these before using them in a
// production environment.
//
//type Middleware struct {
//	logger *logrus.Logger
//}
//
//func (a *Middleware) Logging(next http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		t := time.Now()
//		next(w, r)
//
//		var code int
//		if err := intake.FromContext(r, "response-code", &code); err != nil {
//			a.logger.WithError(err).Error("error getting response code from context")
//		}
//
//		a.logger.WithFields(logrus.Fields{
//			"method":           r.Method,
//			"requestUri":       r.RequestURI,
//			"contentLen":       r.ContentLength,
//			"responseTimeMs":   time.Now().Sub(t).Milliseconds(),
//			"code":             code,
//		}).Info("handled request")
//	}
//}
//
//// Recover avoids the application panicing if any calls to the route cause one
//func (a *Middleware) Recover(next http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		next(w, r)
//		if err := recover(); err != nil {
//			intake.RespondJSON(w, r, http.StatusInternalServerError, "server recovered from a panic")
//			a.logger.WithFields(logrus.Fields{"panic": err}).Error("recovered from panic")
//		}
//	}
//}
//
//// RateLimit will limit requests that use this middleware to n requests per second
//func (a *Middleware) RateLimit(n float64) func(handler http.HandlerFunc) http.HandlerFunc {
//	var lastRequestTime = time.Now()
//	return func(next http.HandlerFunc) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			requestsPerSecond := 1 / time.Now().Sub(lastRequestTime).Seconds()
//			lastRequestTime = time.Now()
//			if requestsPerSecond > n {
//				intake.RespondJSON(w, r, http.StatusTooManyRequests, "rate limited")
//				a.logger.WithFields(logrus.Fields{"requestsPerSecond": requestsPerSecond}).Warn("rate limited")
//				return
//			}
//			next(w, r)
//		}
//	}
//}
//
//// RateLimitIP limit the number of requests per second per IP
//func (a *Middleware) RateLimitIP(n float64) func(handler http.HandlerFunc) http.HandlerFunc {
//	var ipMap map[string]time.Time
//	ipMap = make(map[string]time.Time)
//
//	fn := func(s string) string {
//		ipSplit := strings.Split(s, ":")
//		if len(ipSplit) > 0 {
//			ip := net.ParseIP(ipSplit[0])
//			if ip != nil {
//				return ip.String()
//			}
//		}
//
//		return "noip"
//	}
//	return func(next http.HandlerFunc) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			var requestsPerSecond float64
//			requestIp := fn(r.RemoteAddr)
//			val, exists := ipMap[requestIp]
//			if exists != true {
//				requestsPerSecond = 0
//			} else {
//				requestsPerSecond = 1 / time.Now().Sub(val).Seconds()
//			}
//
//			ipMap[requestIp] = time.Now()
//			if requestsPerSecond > n {
//				intake.RespondJSON(w, r, http.StatusTooManyRequests, "rate limited")
//				a.logger.WithFields(logrus.Fields{"requestsPerSecond": requestsPerSecond, "ip": requestIp}).Warn("ip rate limited")
//				return
//			}
//			next(w, r)
//		}
//	}
//}
//
//// Timeout if added will created a context that will cancel after t. This cancel
//// will affect all downstream uses of the context
//func (a *Middleware) Timeout(t time.Duration) func(handler http.HandlerFunc) http.HandlerFunc {
//	return func(next http.HandlerFunc) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			// Create a context that is both manually cancellable and will signal
//			// cancel at the specified duration.
//			ctx, cancel := context.WithTimeout(r.Context(), t)
//			*r = *r.WithContext(ctx)
//			next(w, r)
//			cancel()
//		}
//	}
//}
