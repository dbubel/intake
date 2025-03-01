package intake

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

	t.Run("test options handler", func(t *testing.T) {
		optionsHandler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Allow", "GET, POST, OPTIONS")
			w.WriteHeader(http.StatusOK)
		}

		app.OptionsHandler(optionsHandler)
		app.AddEndpoint(http.MethodGet, "/test-options", testHandler)
		app.AddEndpoint(http.MethodPost, "/test-options", testHandler)

		r := httptest.NewRequest(http.MethodOptions, "/test-options", nil)
		w := httptest.NewRecorder()
		app.Mux.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		if allow := w.Header().Get("Allow"); allow != "GET, POST, OPTIONS" {
			t.Errorf("Expected Allow header %q, got %q", "GET, POST, OPTIONS", allow)
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

	t.Run("test duplicate options handler", func(t *testing.T) {
		optionsCallCount := 0
		optionsHandler := func(w http.ResponseWriter, r *http.Request) {
			optionsCallCount++
			w.WriteHeader(http.StatusOK)
		}

		app.OptionsHandler(optionsHandler)
		app.AddEndpoint(http.MethodGet, "/test-duplicate", testHandler)
		app.AddEndpoint(http.MethodPost, "/test-duplicate", testHandler)

		r := httptest.NewRequest(http.MethodOptions, "/test-duplicate", nil)
		w := httptest.NewRecorder()
		app.Mux.ServeHTTP(w, r)

		if optionsCallCount != 1 {
			t.Errorf("Expected options handler to be called once, got %d calls", optionsCallCount)
		}
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})
}
