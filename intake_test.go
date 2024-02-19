package intake

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type testPayload struct {
	Msg string `json:"msg"`
}

var (
	payload = testPayload{Msg: "payload"}
	l       *logrus.Logger
)

func init() {
	l = logrus.New()
	l.SetLevel(logrus.InfoLevel)
}

func TestIntake(t *testing.T) {
	app := New(l)

	testHandler := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		RespondJSON(w, r, http.StatusOK, payload)
	}

	app.AddEndpoint(http.MethodGet, "/test", testHandler)

	t.Run("test simple route", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		resp, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var res testPayload
		assert.NoError(t, json.Unmarshal(resp, &res))
		assert.Equal(t, "payload", res.Msg)
    fmt.Println("hello")
     
	})

	t.Run("test route not found", func(t *testing.T) {
    fmt.Println("hello2")
		r := httptest.NewRequest(http.MethodGet, "/not-found", nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestIntakeMiddleware(t *testing.T) {
	app := NewDefault()
	testHandler := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		RespondJSON(w, r, http.StatusOK, payload)
	}
	testHandlerWithCtx := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		var helloStr string
		FromContext(r, "onCtx", &helloStr)
		RespondJSON(w, r, http.StatusOK, helloStr)
	}
	testHandlerWithCtxTwo := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		var helloStr string
		FromContext(r, "onCtx", &helloStr)

		var helloStrTwo string
		FromContext(r, "onCtxTwo", &helloStrTwo)
		RespondJSON(w, r, http.StatusOK, fmt.Sprintf("%s %s", helloStr, helloStrTwo))
	}
	testMw := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			next(w, r, params)
		}
	}
	testMwWithCtx := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			AddToContext(r, "onCtx", "hello world")
			next(w, r, params)
		}
	}
	testMwWithCtxTwo := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			AddToContext(r, "onCtxTwo", "hello world two")
			next(w, r, params)
		}
	}

	app.AddEndpoint(http.MethodGet, "/test-mw-simple", testHandler, testMw)
	app.AddEndpoint(http.MethodGet, "/test-mw-simple-with-context", testHandlerWithCtx, testMwWithCtx)
	app.AddEndpoint(http.MethodGet, "/test-mw-simple-with-context-two", testHandlerWithCtxTwo, testMwWithCtx, testMwWithCtxTwo)

	t.Run("test route with middleware", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/test-mw-simple-with-context", nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("test route with middleware context", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/test-mw-simple-with-context", nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		resp, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var res string
		assert.NoError(t, json.Unmarshal(resp, &res))
		assert.Equal(t, "hello world", res)
	})

	t.Run("test route with middleware with two contexts", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/test-mw-simple-with-context-two", nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		resp, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var res string
		assert.NoError(t, json.Unmarshal(resp, &res))
		assert.Equal(t, "hello world hello world two", res)
	})
}

func TestIntakeMiddlewareGroups(t *testing.T) {
	app := New(l)

	testHandlerWithCtx := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		var helloStr string
		FromContext(r, "onCtx", &helloStr)
		RespondJSON(w, r, http.StatusOK, helloStr)
	}

	testMwWithCtx := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			AddToContext(r, "onCtx", "hello world")
			next(w, r, params)
		}
	}

	eps := Endpoints{
		GET("/test-one", testHandlerWithCtx),
		GET("/test-two", testHandlerWithCtx),
	}
	eps.Use(testMwWithCtx)
	app.AddEndpoints(eps)

	t.Run("test route with middleware group", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/test-one", nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		resp, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var res string
		assert.NoError(t, json.Unmarshal(resp, &res))
		assert.Equal(t, "hello world", res)
	})

	t.Run("test route with middleware group second route", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/test-two", nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		resp, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var res string
		assert.NoError(t, json.Unmarshal(resp, &res))
		assert.Equal(t, "hello world", res)
	})
}

func TestIntakeGlobalMiddleware(t *testing.T) {
	app := New(l)

	testHandlerWithCtx := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		var helloStr string
		FromContext(r, "onCtx", &helloStr)

		var global string
		FromContext(r, "global", &global)
		RespondJSON(w, r, http.StatusOK, fmt.Sprintf("%s %s", helloStr, global))
	}

	testMwWithCtx := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			AddToContext(r, "onCtx", "hello world")
			next(w, r, params)
		}
	}

	testMwGlobal := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			AddToContext(r, "global", "global middleware")
			next(w, r, params)
		}
	}
	app.AddGlobalMiddleware(testMwGlobal)

	eps := Endpoints{
		GET("/test-one", testHandlerWithCtx),
	}
	eps.Use(testMwWithCtx)
	app.AddEndpoints(eps)

	t.Run("test route with middleware group", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/test-one", nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		resp, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		var res string
		assert.NoError(t, json.Unmarshal(resp, &res))
		assert.Equal(t, "hello world global middleware", res)
	})
}
