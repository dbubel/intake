// Package intake provides HTTP routing utilities.
// This file contains middleware implementations for common functionality such as CORS.
package intake

import (
	"net/http"
	"strings"
)

// CORSConfig defines the configuration options for the CORS middleware.
type CORSConfig struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// Default value is ["*"]
	AllowedOrigins []string

	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (GET, POST, HEAD)
	AllowedMethods []string

	// AllowedHeaders is a list of headers the client is allowed to use with
	// cross-domain requests. If the special "*" value is present in the list,
	// all headers will be allowed. Default value is ["Origin", "Accept", "Content-Type"]
	AllowedHeaders []string

	// ExposeHeaders is a list of headers that should be accessible to js in browsers
	ExposeHeaders []string

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached.
	MaxAge int
}

// DefaultCORSConfig returns a default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodHead,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Origin",
			"Accept",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders:    []string{},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours
	}
}

// CORS returns a middleware that handles Cross-Origin Resource Sharing.
// It implements the behavior defined in the W3C CORS specification.
func CORS(config CORSConfig) MiddleWare {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				// Not a CORS request or same origin
				next(w, r)
				return
			}

			// Check if the origin is allowed
			originAllowed := false
			for _, allowedOrigin := range config.AllowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					originAllowed = true
					break
				}
			}

			if !originAllowed {
				next(w, r)
				return
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				// Set standard CORS headers
				corsHeaders(w, config, origin)

				// Handle preflight specific headers
				w.Header().Set("Access-Control-Max-Age", string(config.MaxAge))

				// Set allowed methods
				if len(config.AllowedMethods) > 0 {
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
				}

				// Set allowed headers
				if len(config.AllowedHeaders) > 0 {
					if config.AllowedHeaders[0] == "*" {
						// If we allow all headers, mirror the request headers
						reqHeaders := r.Header.Get("Access-Control-Request-Headers")
						if reqHeaders != "" {
							w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
						}
					} else {
						w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
					}
				}

				// Respond with 204 status and return immediately
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Handle actual request (not a preflight)
			corsHeaders(w, config, origin)
			next(w, r)
		}
	}
}

// corsHeaders sets the common CORS headers on the response
func corsHeaders(w http.ResponseWriter, config CORSConfig, origin string) {
	// Set allowed origin
	if config.AllowedOrigins[0] == "*" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	// Set credentials flag if needed
	if config.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// Set exposed headers if any
	if len(config.ExposeHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
	}
}