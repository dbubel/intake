package intake

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSPreflightHeaderValidation(t *testing.T) {
	app := New()
	app.AddGlobalMiddleware(CORS(CORSConfig{
		AllowedOrigins: []string{"https://example.com"},
		AllowedMethods: []string{http.MethodGet},
		AllowedHeaders: []string{"X-Token", "Content-Type"},
	}))
	app.AddEndpoint(http.MethodGet, "/data", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	app.AddOptionsEndpoints()

	t.Run("rejects disallowed headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/data", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", http.MethodGet)
		req.Header.Set("Access-Control-Request-Headers", "X-Token, X-Other")
		rr := httptest.NewRecorder()

		app.Mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
		}
	})

	t.Run("allows configured headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/data", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", http.MethodGet)
		req.Header.Set("Access-Control-Request-Headers", "x-token")
		rr := httptest.NewRecorder()

		app.Mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
		}
		if got := rr.Header().Get("Access-Control-Allow-Headers"); got == "" {
			t.Fatalf("expected Access-Control-Allow-Headers, got empty")
		}
	})
}

func TestCORSWildcardHeadersEcho(t *testing.T) {
	app := New()
	app.AddGlobalMiddleware(CORS(CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet},
		AllowedHeaders: []string{"*"},
	}))
	app.AddEndpoint(http.MethodGet, "/data", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	app.AddOptionsEndpoints()

	req := httptest.NewRequest(http.MethodOptions, "/data", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	req.Header.Set("Access-Control-Request-Headers", "X-Foo, X-Bar")
	rr := httptest.NewRecorder()

	app.Mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
	if got := rr.Header().Get("Access-Control-Allow-Headers"); got != "X-Foo, X-Bar" {
		t.Fatalf("expected Access-Control-Allow-Headers to echo request, got %q", got)
	}
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("expected Access-Control-Allow-Origin '*', got %q", got)
	}
}

func TestCORSWildcardOriginSchemeMatch(t *testing.T) {
	app := New()
	app.AddGlobalMiddleware(CORS(CORSConfig{
		AllowedOrigins: []string{"http://*.example.com"},
		AllowedMethods: []string{http.MethodGet},
	}))
	app.AddEndpoint(http.MethodGet, "/data", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	app.AddOptionsEndpoints()

	t.Run("allows http subdomain", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/data", nil)
		req.Header.Set("Origin", "http://api.example.com")
		rr := httptest.NewRecorder()

		app.Mux.ServeHTTP(rr, req)

		if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "http://api.example.com" {
			t.Fatalf("expected Access-Control-Allow-Origin to echo origin, got %q", got)
		}
	})

	t.Run("rejects https subdomain", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/data", nil)
		req.Header.Set("Origin", "https://api.example.com")
		rr := httptest.NewRecorder()

		app.Mux.ServeHTTP(rr, req)

		if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Fatalf("expected no Access-Control-Allow-Origin, got %q", got)
		}
	})
}
