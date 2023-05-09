package intake

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {

	router := NewRouter()
	router.AddRoute("/api/v1/foo", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "foo")
	})
	router.AddRoute("/api/v1/bar", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "bar")
	})
	router.AddRoute("/api/v1/baz/bar", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "bar/bar")
	})
	router.AddRoute("/world", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("World!"))
	})
	router.AddRoute("/hello", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello!"))
	})

	t.Run("test hello route", func(t *testing.T) {
		server := httptest.NewServer(router)
		defer server.Close()
		resp, err := http.Get(server.URL + "/hello")
		body, err := io.ReadAll(resp.Body)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, err)
		assert.Equal(t, "Hello!", string(body))
	})

	t.Run("test world route", func(t *testing.T) {
		server := httptest.NewServer(router)
		defer server.Close()
		resp, err := http.Get(server.URL + "/world")
		body, err := io.ReadAll(resp.Body)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, err)
		assert.Equal(t, "World!", string(body))
	})

	t.Run("test api/foo route", func(t *testing.T) {
		server := httptest.NewServer(router)
		defer server.Close()
		resp, err := http.Get(server.URL + "/api/v1/foo")
		body, err := io.ReadAll(resp.Body)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, err)
		assert.Equal(t, "foo", string(body))
	})

	t.Run("test 404 not found", func(t *testing.T) {
		server := httptest.NewServer(router)
		defer server.Close()
		resp, err := http.Get(server.URL + "/api/v1/does_not_exist")
		body, err := io.ReadAll(resp.Body)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.NoError(t, err)
		assert.Equal(t, "404 page not found\n", string(body))
	})
}
