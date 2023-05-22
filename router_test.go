package intake

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

	t.Run("test 301 redirect trailing slash no follow", func(t *testing.T) {
		server := httptest.NewServer(router)
		defer server.Close()
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		resp, err := client.Get(server.URL + "/hello/")
		if resp.StatusCode != http.StatusMovedPermanently {
			t.Errorf("status code != 301 [%d]", resp.StatusCode)
		}

		if err != nil {
			t.Error(err)
		}
	})

	t.Run("test follow redirect with slash", func(t *testing.T) {
		server := httptest.NewServer(router)
		defer server.Close()
		resp, err := http.Get(server.URL + "/hello/")
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
}

func TestAddRoute(t *testing.T) {
	router := NewRouter()
	route := "/test"
	method := http.MethodGet
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {}
	router.AddRoute(route, method, handler)

	// Assert that the route was added to the router
	if _, found := router.routes[method+strings.TrimSuffix(route, "/")]; !found {
		t.Errorf("Route was not added to the router")
	}
}

func TestServeHTTP(t *testing.T) {
	router := NewRouter()
	route := "/test"
	method := http.MethodGet
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	router.AddRoute(route, method, handler)

	req, err := http.NewRequest(method, route, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler = router.routes[method+route]

	handler(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestRedirectTrailingSlash(t *testing.T) {
	//router := NewRouter()

	req, err := http.NewRequest(http.MethodGet, "/test/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	redirect := RedirectTrailingSlash(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("RedirectTrailingSlash returned wrong status code: got %v want %v", status, http.StatusMovedPermanently)
	}

	if !redirect {
		t.Errorf("Expected a redirect but did not get one")
	}
}
