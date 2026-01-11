package intake

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkServeHTTP(b *testing.B) {
	app := New()
	app.AddGlobalMiddleware(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r)
		}
	})
	app.AddEndpoint(http.MethodGet, "/bench", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/bench", nil)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		app.Mux.ServeHTTP(rr, req)
	}
}

func BenchmarkCORSPreflight(b *testing.B) {
	app := New()
	app.AddGlobalMiddleware(CORS(CORSConfig{
		AllowedOrigins: []string{"https://example.com"},
		AllowedMethods: []string{http.MethodGet, http.MethodOptions},
		AllowedHeaders: []string{"X-Token", "Content-Type"},
		MaxAge:         600,
	}))
	app.AddEndpoint(http.MethodGet, "/bench", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodOptions, "/bench", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	req.Header.Set("Access-Control-Request-Headers", "X-Token")
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		app.Mux.ServeHTTP(rr, req)
	}
}
