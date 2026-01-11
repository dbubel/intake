// Package intake provides HTTP routing utilities.
// This file contains middleware implementations for Cross-Origin Resource Sharing (CORS).
package intake

import (
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

type originPattern struct {
	scheme string
	suffix string
}

type corsPolicy struct {
	allowedMethods     []string
	allowedMethodsSet  map[string]struct{}
	allowedMethodsHeader string
	allowedHeaders     []string
	allowedHeadersSet  map[string]struct{}
	allowedHeadersHeader string
	allowedOrigins     map[string]struct{}
	allowedPatterns    []originPattern
	allowAnyOrigin     bool
	allowAnyHeader     bool
	exposeHeaders      []string
	exposeHeadersHeader string
	allowCredentials   bool
	maxAge             int
}

// CORSConfig defines the configuration options for the CORS middleware.
// This struct allows for fine-grained control over CORS policy implementation.
// Each field corresponds to a specific CORS header or behavior as defined in
// the W3C CORS specification (https://www.w3.org/TR/cors/).
type CORSConfig struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// Default value is ["*"], which allows any origin.
	// Examples: ["https://example.com", "https://*.example.com", "*"]
	AllowedOrigins []string

	// AllowedMethods is a list of HTTP methods the client is allowed to use with
	// cross-domain requests. Default value includes all simple methods (GET, POST, HEAD)
	// plus commonly used methods (PUT, DELETE, PATCH, OPTIONS).
	// This controls the Access-Control-Allow-Methods header.
	AllowedMethods []string

	// AllowedHeaders is a list of headers the client is allowed to use with
	// cross-domain requests. If the special "*" value is present in the list,
	// all headers will be allowed. Default value is ["Origin", "Accept", "Content-Type", "Authorization"].
	// This controls the Access-Control-Allow-Headers header.
	AllowedHeaders []string

	// ExposeHeaders is a list of headers that should be accessible to JavaScript in browsers.
	// These headers will be listed in the Access-Control-Expose-Headers response header.
	// By default, no headers are exposed.
	ExposeHeaders []string

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	// This controls the Access-Control-Allow-Credentials header.
	// Note: Cannot be used with wildcard (*) AllowedOrigins for security reasons.
	AllowCredentials bool

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached by the browser. Default is 86400 seconds (24 hours).
	// This controls the Access-Control-Max-Age header.
	MaxAge int
}

// DefaultCORSConfig returns a default CORS configuration with common settings.
// The default configuration:
// - Allows all origins (*)
// - Includes all standard HTTP methods
// - Sets commonly used headers
// - Disables credentials
// - Sets a 24-hour cache period for preflight requests
//
// This provides a secure starting point that can be customized as needed.
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
// It implements the behavior defined in the W3C CORS specification (https://www.w3.org/TR/cors/).
//
// This middleware handles both preflight OPTIONS requests and actual CORS requests:
// - For preflight requests, it validates requested methods and headers
// - For actual requests, it applies appropriate CORS headers based on configuration
// - It supports wildcard origins, domain pattern matching, and specific origin lists
// - It ensures compliance with security requirements (e.g., disallowing credentials with wildcard origins)
//
// Parameters:
//   - config: The CORSConfig struct containing CORS policy configuration
//
// Returns:
//   - A MiddleWare function that can be applied to HTTP handlers
func CORS(config CORSConfig) MiddleWare {
	// Validate the configuration
	// Ensure we have at least one allowed method if not explicitly set
	if len(config.AllowedMethods) == 0 {
		config.AllowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodHead}
	}

	policy := buildPolicy(config)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				// Not a CORS request or same origin request - proceed without CORS headers
				next(w, r)
				return
			}

			// Check if the origin is allowed by the configured policy
			originAllowed := policy.isOriginAllowed(origin)
			if !originAllowed {
				// Origin not allowed, pass through without CORS headers
				// This maintains security by not acknowledging invalid cross-origin requests
				next(w, r)
				return
			}

			// Handle preflight OPTIONS requests
			// Preflight requests are sent by browsers before the actual request to check
			// if the CORS request is allowed by the server
			if r.Method == http.MethodOptions {
				// Set standard CORS headers for all responses
				corsHeaders(w, policy, origin)

				// Set cache duration for preflight response
				// This helps reduce the number of preflight requests
				if policy.maxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", strconv.Itoa(policy.maxAge))
				}

				// Check if the requested HTTP method is allowed
				requestMethod := r.Header.Get("Access-Control-Request-Method")
				if requestMethod != "" {
					_, methodAllowed := policy.allowedMethodsSet[requestMethod]
					if !methodAllowed {
						// Method not allowed - respond with 403 Forbidden
						w.WriteHeader(http.StatusForbidden)
						return
					}
				}

				// Set the list of allowed HTTP methods
				if policy.allowedMethodsHeader != "" {
					w.Header().Set("Access-Control-Allow-Methods", policy.allowedMethodsHeader)
				}

				// Handle the requested headers check
				requestHeaders := r.Header.Get("Access-Control-Request-Headers")
				if len(policy.allowedHeaders) > 0 || policy.allowAnyHeader {
					if policy.allowAnyHeader {
						// If wildcard is configured for headers, mirror the requested headers
						// This allows the browser to send any headers it needs
						if requestHeaders != "" {
							w.Header().Set("Access-Control-Allow-Headers", requestHeaders)
						}
					} else {
						// Otherwise, only allow the specifically configured headers,
						// and reject preflights that ask for disallowed headers.
						if requestHeaders != "" && !policy.areHeadersAllowed(requestHeaders) {
							w.WriteHeader(http.StatusForbidden)
							return
						}
						if policy.allowedHeadersHeader != "" {
							w.Header().Set("Access-Control-Allow-Headers", policy.allowedHeadersHeader)
						}
					}
				} else if requestHeaders != "" {
					// No allowed headers configured: reject explicit header requests.
					w.WriteHeader(http.StatusForbidden)
					return
				}

				// Preflight requests only need headers, not content
				// Respond with 204 No Content status and return immediately
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Handle actual CORS request (not a preflight)
			// Apply the CORS headers and continue with request processing
			corsHeaders(w, policy, origin)
			next(w, r)
		}
	}
}

// corsHeaders sets the common CORS headers on the response.
// This internal helper function is used to consistently apply the basic
// CORS headers required for both preflight and actual CORS requests.
//
// Parameters:
//   - w: The HTTP response writer to set headers on
//   - config: The CORS policy to apply
//   - origin: The requesting Origin header value
func corsHeaders(w http.ResponseWriter, config corsPolicy, origin string) {
	// Set Access-Control-Allow-Origin header
	// There are two strategies based on configuration:
	// 1. Use "*" when wildcard origins are allowed and credentials aren't required
	// 2. Mirror the specific origin otherwise (required when using credentials)
	if config.allowAnyOrigin && !config.allowCredentials {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		// Echo back the specific origin
		w.Header().Set("Access-Control-Allow-Origin", origin)
		// When returning a specific origin, Vary header is required for proper caching
		// This prevents cache poisoning across different origins
		w.Header().Add("Vary", "Origin")
	}

	// Set Access-Control-Allow-Credentials header if credentials are allowed
	// This enables sending cookies, authorization headers, and TLS client certs
	if config.allowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// Set Access-Control-Expose-Headers header if specific headers should be
	// accessible to JavaScript in the browser
	if config.exposeHeadersHeader != "" {
		w.Header().Set("Access-Control-Expose-Headers", config.exposeHeadersHeader)
	}
}

func buildPolicy(config CORSConfig) corsPolicy {
	policy := corsPolicy{
		allowedMethods:    config.AllowedMethods,
		allowedMethodsSet: make(map[string]struct{}, len(config.AllowedMethods)),
		allowedMethodsHeader: strings.Join(config.AllowedMethods, ", "),
		allowedHeaders:    config.AllowedHeaders,
		allowedHeadersSet: make(map[string]struct{}, len(config.AllowedHeaders)),
		allowedHeadersHeader: strings.Join(config.AllowedHeaders, ", "),
		allowedOrigins:    make(map[string]struct{}, len(config.AllowedOrigins)),
		allowAnyOrigin:    false,
		allowAnyHeader:    containsWildcard(config.AllowedHeaders),
		exposeHeaders:     config.ExposeHeaders,
		exposeHeadersHeader: strings.Join(config.ExposeHeaders, ", "),
		allowCredentials:  config.AllowCredentials,
		maxAge:            config.MaxAge,
	}

	for _, method := range config.AllowedMethods {
		policy.allowedMethodsSet[method] = struct{}{}
	}

	for _, header := range config.AllowedHeaders {
		policy.allowedHeadersSet[strings.ToLower(header)] = struct{}{}
	}

	for _, origin := range config.AllowedOrigins {
		if origin == "*" {
			policy.allowAnyOrigin = true
			continue
		}

		if strings.HasPrefix(origin, "https://*.") {
			policy.allowedPatterns = append(policy.allowedPatterns, originPattern{
				scheme: "https",
				suffix: origin[len("https://*."):],
			})
			continue
		}
		if strings.HasPrefix(origin, "http://*.") {
			policy.allowedPatterns = append(policy.allowedPatterns, originPattern{
				scheme: "http",
				suffix: origin[len("http://*."):],
			})
			continue
		}

		policy.allowedOrigins[origin] = struct{}{}
	}

	// Invalid configuration: wildcard origin with credentials.
	// Remove wildcard to maintain security.
	if policy.allowCredentials && policy.allowAnyOrigin {
		policy.allowAnyOrigin = false
	}

	return policy
}

func (p corsPolicy) isOriginAllowed(origin string) bool {
	if p.allowAnyOrigin {
		return true
	}
	if _, ok := p.allowedOrigins[origin]; ok {
		return true
	}
	if len(p.allowedPatterns) == 0 {
		return false
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	host := u.Hostname()
	if host == "" {
		return false
	}

	for _, pattern := range p.allowedPatterns {
		if u.Scheme != pattern.scheme {
			continue
		}
		if host != pattern.suffix && strings.HasSuffix(host, "."+pattern.suffix) {
			return true
		}
	}
	return false
}

func (p corsPolicy) areHeadersAllowed(requestHeaders string) bool {
	if requestHeaders == "" {
		return true
	}
	for _, h := range strings.Split(requestHeaders, ",") {
		header := strings.ToLower(strings.TrimSpace(h))
		if header == "" {
			continue
		}
		if _, ok := p.allowedHeadersSet[header]; !ok {
			return false
		}
	}
	return true
}

// containsWildcard checks if the slice contains the wildcard "*" value.
// This is a helper function used to determine if wildcard patterns exist
// in configuration settings like AllowedOrigins or AllowedHeaders.
//
// Parameters:
//   - s: The string slice to check for wildcards
//
// Returns:
//   - true if the slice contains "*", false otherwise
func containsWildcard(s []string) bool {
	return slices.Contains(s, "*")
}
