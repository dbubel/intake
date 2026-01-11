package intake

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testPayload struct {
	Msg string `json:"msg"`
}

func TestIntake(t *testing.T) {
	payload := testPayload{Msg: "test response"}
	var app = New()

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(payload)
	}

	t.Run("test single endpoint", func(t *testing.T) {
		app.AddEndpoint(http.MethodGet, "/test", testHandler)

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		app.Mux.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		resp, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var res testPayload
		if err := json.Unmarshal(resp, &res); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if res.Msg != "test response" {
			t.Errorf("Expected message %q, got %q", "test response", res.Msg)
		}
	})

	t.Run("test path params", func(t *testing.T) {
		paramApp := New()
		var got string
		handler := func(w http.ResponseWriter, r *http.Request) {
			got = r.PathValue("hello")
			w.WriteHeader(http.StatusOK)
		}

		paramApp.AddEndpoint(http.MethodGet, "/api/{hello}/world", handler)

		r := httptest.NewRequest(http.MethodGet, "/api/hi/world", nil)
		w := httptest.NewRecorder()
		paramApp.Mux.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		if got != "hi" {
			t.Errorf("Expected path value %q, got %q", "hi", got)
		}
	})

	t.Run("test middleware execution", func(t *testing.T) {
		middlewareCalled := false
		middleware := func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				middlewareCalled = true
				next(w, r)
			}
		}

		app.AddGlobalMiddleware(middleware)
		app.AddEndpoint(http.MethodGet, "/test-middleware", testHandler)

		r := httptest.NewRequest(http.MethodGet, "/test-middleware", nil)
		w := httptest.NewRecorder()
		app.Mux.ServeHTTP(w, r)

		if !middlewareCalled {
			t.Error("Expected middleware to be called, but it wasn't")
		}
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("test adding multiple endpoints", func(t *testing.T) {
		handler1Called := false
		handler2Called := false

		handler1 := func(w http.ResponseWriter, r *http.Request) {
			handler1Called = true
			w.WriteHeader(http.StatusOK)
		}

		handler2 := func(w http.ResponseWriter, r *http.Request) {
			handler2Called = true
			w.WriteHeader(http.StatusCreated)
		}

		endpoints := Endpoints{
			{
				Verb:            http.MethodGet,
				Path:            "/multiple1",
				EndpointHandler: handler1,
			},
			{
				Verb:            http.MethodPost,
				Path:            "/multiple2",
				EndpointHandler: handler2,
			},
		}

		app.AddEndpoints(endpoints)

		// Test first endpoint
		r1 := httptest.NewRequest(http.MethodGet, "/multiple1", nil)
		w1 := httptest.NewRecorder()
		app.Mux.ServeHTTP(w1, r1)

		if !handler1Called {
			t.Error("Expected handler1 to be called, but it wasn't")
		}
		if w1.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w1.Code)
		}

		// Test second endpoint
		r2 := httptest.NewRequest(http.MethodPost, "/multiple2", nil)
		w2 := httptest.NewRecorder()
		app.Mux.ServeHTTP(w2, r2)

		if !handler2Called {
			t.Error("Expected handler2 to be called, but it wasn't")
		}
		if w2.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w2.Code)
		}
	})

	t.Run("test CORS middleware", func(t *testing.T) {
		newApp := New()

		// Create CORS middleware
		corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodOptions {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
					w.WriteHeader(http.StatusOK)
					return
				}

				// Set CORS headers for non-OPTIONS methods as well
				w.Header().Set("Access-Control-Allow-Origin", "*")
				next(w, r)
			}
		}

		// Apply CORS middleware globally
		newApp.AddGlobalMiddleware(corsMiddleware)

		// Add test endpoints with explicit support for OPTIONS method
		corsHandler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}
		newApp.AddEndpoint(http.MethodGet, "/cors-test", corsHandler)
		newApp.AddEndpoint(http.MethodOptions, "/cors-test", corsHandler)

		// Test OPTIONS request
		optionsReq := httptest.NewRequest(http.MethodOptions, "/cors-test", nil)
		optionsResp := httptest.NewRecorder()
		newApp.Mux.ServeHTTP(optionsResp, optionsReq)

		// Verify response code for OPTIONS
		if optionsResp.Code != http.StatusOK {
			t.Errorf("Expected status code %d for OPTIONS request, got %d", http.StatusOK, optionsResp.Code)
		}

		// Verify CORS headers for OPTIONS
		if optionsResp.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("Expected Access-Control-Allow-Origin header to be '*', got '%s'",
				optionsResp.Header().Get("Access-Control-Allow-Origin"))
		}

		if !strings.Contains(optionsResp.Header().Get("Access-Control-Allow-Methods"), "OPTIONS") {
			t.Errorf("Expected Access-Control-Allow-Methods to contain 'OPTIONS', got '%s'",
				optionsResp.Header().Get("Access-Control-Allow-Methods"))
		}

		// Test regular GET request
		getReq := httptest.NewRequest(http.MethodGet, "/cors-test", nil)
		getResp := httptest.NewRecorder()
		newApp.Mux.ServeHTTP(getResp, getReq)

		// Verify response for GET
		if getResp.Code != http.StatusOK {
			t.Errorf("Expected status code %d for GET request, got %d", http.StatusOK, getResp.Code)
		}

		// Verify CORS headers are also present on non-OPTIONS responses
		if getResp.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("Expected Access-Control-Allow-Origin header to be '*' for GET request, got '%s'",
				getResp.Header().Get("Access-Control-Allow-Origin"))
		}

		// Verify response body
		body, _ := ioutil.ReadAll(getResp.Body)
		if string(body) != "success" {
			t.Errorf("Expected body 'success', got '%s'", string(body))
		}
	})

	t.Run("test panic handler", func(t *testing.T) {
		panicApp := New()
		panicHandlerCalled := false
		errorMessage := "Something went wrong"

		// Set up a panic handler
		panicApp.SetPanicHandler(func(w http.ResponseWriter, r *http.Request, err any) {
			panicHandlerCalled = true
			w.WriteHeader(http.StatusInternalServerError)
			errMsg, ok := err.(string)
			if !ok {
				errMsg = "Unknown error"
			}
			w.Write([]byte(errMsg))
		})

		// Create a handler that will panic
		panicHandler := func(w http.ResponseWriter, r *http.Request) {
			panic(errorMessage)
		}

		// Register the panic-causing endpoint
		panicApp.AddEndpoint(http.MethodGet, "/panic", panicHandler)

		// Make a request to the endpoint
		r := httptest.NewRequest(http.MethodGet, "/panic", nil)
		w := httptest.NewRecorder()
		panicApp.Mux.ServeHTTP(w, r)

		// Verify the panic was caught
		if !panicHandlerCalled {
			t.Error("Expected panic handler to be called, but it wasn't")
		}

		// Verify response code
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}

		// Verify response body
		body, _ := ioutil.ReadAll(w.Body)
		if string(body) != errorMessage {
			t.Errorf("Expected body '%s', got '%s'", errorMessage, string(body))
		}
	})

	t.Run("test panic handler with middleware", func(t *testing.T) {
		panicApp := New()
		middlewareCalled := false
		panicHandlerCalled := false
		errorMessage := "Middleware panic"

		// Set up a panic handler
		panicApp.SetPanicHandler(func(w http.ResponseWriter, r *http.Request, err any) {
			panicHandlerCalled = true

			w.WriteHeader(http.StatusInternalServerError)
			errMsg, ok := err.(string)
			if !ok {
				errMsg = "Unknown error"
			}
			w.Write([]byte(errMsg))
		})

		// Create middleware that marks execution
		requestIDMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				middlewareCalled = true

				next(w, r)
			}
		}

		// Create another middleware that will panic
		panicMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				// This will panic before reaching the handler
				panic(errorMessage)
			}
		}

		// Add the middleware to the app
		panicApp.AddGlobalMiddleware(requestIDMiddleware)

		// Create a normal handler
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("This should not be called"))
		}

		// Register the endpoint with the panic middleware
		panicApp.AddEndpoint(http.MethodGet, "/middleware-panic", handler, panicMiddleware)

		// Make a request to the endpoint
		r := httptest.NewRequest(http.MethodGet, "/middleware-panic", nil)
		w := httptest.NewRecorder()
		panicApp.Mux.ServeHTTP(w, r)

		// Verify the middleware was called
		if !middlewareCalled {
			t.Error("Expected middleware to be called, but it wasn't")
		}

		// Verify the panic was caught
		if !panicHandlerCalled {
			t.Error("Expected panic handler to be called, but it wasn't")
		}

		// Verify response code
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}

		// Verify response body
		body, _ := ioutil.ReadAll(w.Body)
		if string(body) != errorMessage {
			t.Errorf("Expected body '%s', got '%s'", errorMessage, string(body))
		}
	})
}

func TestIsOriginAllowedWildcardSubdomain(t *testing.T) {
	policy := buildPolicy(CORSConfig{
		AllowedOrigins: []string{"https://*.example.com"},
	})

	cases := []struct {
		origin string
		want   bool
	}{
		{origin: "https://api.example.com", want: true},
		{origin: "https://example.com", want: false},
		{origin: "https://badexample.com", want: false},
		{origin: "http://api.example.com", want: false},
	}

	for _, tc := range cases {
		if got := policy.isOriginAllowed(tc.origin); got != tc.want {
			t.Fatalf("origin %q allowed=%v, want %v", tc.origin, got, tc.want)
		}
	}
}
