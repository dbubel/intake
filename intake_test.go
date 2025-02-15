package intake

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
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

		require.Equal(t, http.StatusOK, w.Code)
		resp, err := ioutil.ReadAll(w.Body)
		require.NoError(t, err)

		var res testPayload
		require.NoError(t, json.Unmarshal(resp, &res))
		require.Equal(t, "test response", res.Msg)
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

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "GET, POST, OPTIONS", w.Header().Get("Allow"))
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

		require.True(t, middlewareCalled)
		require.Equal(t, http.StatusOK, w.Code)
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

		require.Equal(t, 1, optionsCallCount)
		require.Equal(t, http.StatusOK, w.Code)
	})
}
