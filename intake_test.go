package intake

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntake_EndpointGroups(t *testing.T) {
	intake := New()

	grpOne := Endpoints{
		GET("/get", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprint(writer, "GET()")
		}),
	}
	grpTwo := Endpoints{
		POST("/post", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprint(writer, "POST()")
		}),
	}
	grpThree := Endpoints{
		DELETE("/delete", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprint(writer, "DELETE()")
		}),
	}

	intake.AddEndpoints(grpOne)
	intake.AddEndpoints(grpTwo)
	intake.AddEndpoints(grpThree)

	server := httptest.NewServer(intake.Router)
	defer server.Close()

	tests := []struct {
		name   string
		resp   string
		route  string
		method string
	}{
		{
			name:   "GET group",
			resp:   "GET()",
			route:  "/get",
			method: http.MethodGet,
		},
		{
			name:   "POST group",
			resp:   "POST()",
			route:  "/post",
			method: http.MethodPost,
		},
		{
			name:   "DELETE group",
			resp:   "DELETE()",
			route:  "/delete",
			method: http.MethodDelete,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, server.URL+test.route, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if http.StatusOK != resp.StatusCode {
				t.Error("incorrect status code", resp.StatusCode)
			}

			if test.resp != string(body) {
				t.Error("incorrect response body", string(body), "expected", test.resp)
			}
		})
	}
}

func TestIntake_HttpMethodWrappers(t *testing.T) {
	intake := New()

	eps := Endpoints{
		GET("/get", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprint(writer, "GET()")
		}),
		POST("/post", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprint(writer, "POST()")
		}),
		PATCH("/patch", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprint(writer, "PATCH()")
		}),
		PUT("/put", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprint(writer, "PUT()")
		}),
		DELETE("/delete", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprint(writer, "DELETE()")
		}),
	}
	intake.AddEndpoints(eps)

	server := httptest.NewServer(intake.Router)
	defer server.Close()

	tests := []struct {
		name   string
		resp   string
		route  string
		method string
	}{
		{
			name:   "GET wrapper",
			resp:   "GET()",
			route:  "/get",
			method: http.MethodGet,
		},
		{
			name:   "POST wrapper",
			resp:   "POST()",
			route:  "/post",
			method: http.MethodPost,
		},
		{
			name:   "PUT wrapper",
			resp:   "PUT()",
			route:  "/put",
			method: http.MethodPut,
		},
		{
			name:   "PATCH wrapper",
			resp:   "PATCH()",
			route:  "/patch",
			method: http.MethodPatch,
		},
		{
			name:   "DELETE wrapper",
			resp:   "DELETE()",
			route:  "/delete",
			method: http.MethodDelete,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, server.URL+test.route, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if http.StatusOK != resp.StatusCode {
				t.Error("incorrect status code", resp.StatusCode)
			}

			if test.resp != string(body) {
				t.Error("incorrect response body", string(body), "expected", test.resp)
			}
		})
	}
}

func TestIntake_Methods(t *testing.T) {
	intake := New()

	intake.AddEndpoint("/get", http.MethodGet, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "hello get")
	})
	intake.AddEndpoint("/post", http.MethodPost, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "hello post")
	})
	intake.AddEndpoint("/patch", http.MethodPatch, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "hello patch")
	})
	intake.AddEndpoint("/put", http.MethodPut, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "hello put")
	})
	intake.AddEndpoint("/delete", http.MethodDelete, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "hello delete")
	})

	server := httptest.NewServer(intake.Router)
	defer server.Close()

	t.Run("test method get", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/get")
		if err != nil {
			t.Error(err)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("status code != 200 [%d]", resp.StatusCode)
		}
		if err != nil {
			t.Error(err)
		}

		if "hello get" != string(body) {
			t.Errorf("body != 'hello get' %s", string(body))
		}
	})

	t.Run("test method post", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPost, server.URL+"/post", nil)
		if err != nil {
			t.Error(err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		body, err := io.ReadAll(resp.Body)
		respBody := string(body)
		if err != nil {
			t.Error(err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Error(err)
			return
		}
		correctResp := "hello post"
		if correctResp != respBody {
			t.Errorf("[%s] != [%s]", correctResp, respBody)
			return
		}
	})

	t.Run("test method patch", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPatch, server.URL+"/patch", nil)
		if err != nil {
			t.Error(err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		body, err := io.ReadAll(resp.Body)
		respBody := string(body)
		if err != nil {
			t.Error(err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Error(err)
			return
		}
		correctResp := "hello patch"
		if correctResp != respBody {
			t.Errorf("[%s] != [%s]", correctResp, respBody)
			return
		}
	})

	t.Run("test method put", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, server.URL+"/put", nil)
		if err != nil {
			t.Error(err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		body, err := io.ReadAll(resp.Body)
		respBody := string(body)
		if err != nil {
			t.Error(err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Error(err)
			return
		}
		correctResp := "hello put"
		if correctResp != respBody {
			t.Errorf("[%s] != [%s]", correctResp, respBody)
			return
		}
	})

	t.Run("test method delete", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, server.URL+"/delete", nil)
		if err != nil {
			t.Error(err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}

		body, err := io.ReadAll(resp.Body)
		respBody := string(body)
		if err != nil {
			t.Error(err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Error(err)
			return
		}
		correctResp := "hello delete"
		if correctResp != respBody {
			t.Errorf("[%s] != [%s]", correctResp, respBody)
			return
		}
	})
}

func TestIntake_AddEndpoint(t *testing.T) {
	intake := New()
	simpleHandler := func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "hello world")
	}

	simpleMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "hello middleware ")
			next(w, r)
		}
	}

	thirdMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "hello middleware three ")
			next(w, r)
		}
	}

	intake.AddEndpoint("/test", http.MethodGet, simpleHandler)
	intake.AddEndpoint("/middleware-simple", http.MethodGet, simpleHandler, simpleMiddleware)
	intake.AddEndpoint("/middleware-simple-three", http.MethodGet, simpleHandler, simpleMiddleware, thirdMiddleware)

	server := httptest.NewServer(intake.Router)
	defer server.Close()

	t.Run("test simple route", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/test")
		body, err := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("status code != 200 [%d]", resp.StatusCode)
		}

		if err != nil {
			t.Error(err)
		}

		if "hello world" != string(body) {
			t.Errorf("body != 'hello world' %s", string(body))
		}
	})

	t.Run("test a single simple middleware", func(t *testing.T) {

		resp, err := http.Get(server.URL + "/middleware-simple")
		body, err := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("status code != 200 [%d]", resp.StatusCode)
		}

		if err != nil {
			t.Error(err)
		}

		if "hello middleware hello world" != string(body) {
			t.Errorf("body != 'hello middleware hello world' %s", string(body))
		}
	})

	t.Run("test multiple simple middleware", func(t *testing.T) {

		resp, err := http.Get(server.URL + "/middleware-simple-three")
		body, err := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("status code != 200 [%d]", resp.StatusCode)
		}

		if err != nil {
			t.Error(err)
		}

		if "hello middleware hello middleware three hello world" != string(body) {
			t.Errorf("body != 'hello middleware hello middleware three hello world' %s", string(body))
		}
	})
}
