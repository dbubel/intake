package intake

import (
	"fmt"
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
		respBody := string(body)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("status code != 200 [%d]", resp.StatusCode)
		}

		if err != nil {
			t.Error(err)
		}

		correctResp := "Hello!"
		if correctResp != respBody {
			t.Errorf("[%s] != [%s]", correctResp, respBody)
		}
	})

	t.Run("test world route", func(t *testing.T) {
		server := httptest.NewServer(router)
		defer server.Close()
		resp, err := http.Get(server.URL + "/world")
		body, err := io.ReadAll(resp.Body)
		respBody := string(body)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("status code != 200 [%d]", resp.StatusCode)
		}

		if err != nil {
			t.Error(err)
		}

		correctResp := "World!"
		if correctResp != respBody {
			t.Errorf("[%s] != [%s]", correctResp, respBody)
		}
	})

	t.Run("test api/foo route", func(t *testing.T) {
		server := httptest.NewServer(router)
		defer server.Close()
		resp, err := http.Get(server.URL + "/api/v1/foo")
		body, err := io.ReadAll(resp.Body)
		respBody := string(body)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("status code != 200 [%d]", resp.StatusCode)
		}

		if err != nil {
			t.Error(err)
		}

		correctResp := "foo"
		if correctResp != respBody {
			t.Errorf("[%s] != [%s]", correctResp, respBody)
		}
	})

	t.Run("test 404 not found", func(t *testing.T) {
		server := httptest.NewServer(router)
		defer server.Close()
		resp, err := http.Get(server.URL + "/api/v1/does_not_exist")
		body, err := io.ReadAll(resp.Body)
		respBody := string(body)

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("status code != 200 [%d]", resp.StatusCode)
		}

		if err != nil {
			t.Error(err)
		}

		correctResp := "404 page not found\n"
		if correctResp != respBody {
			t.Errorf("[%s] != [%s]", correctResp, respBody)
		}
	})
}
